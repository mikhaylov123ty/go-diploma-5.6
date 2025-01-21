package accrual

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/balance"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/orders"

	"github.com/go-resty/resty/v2"
)

const (
	accrualRoot = "/api/orders/"
)

type Accrual struct {
	address     string
	ordersRepo  orders.Storage
	balanceRepo balance.Storage
}

func NewAccrual(address string, ordersRepo orders.Storage, balanceRepo balance.Storage) *Accrual {
	return &Accrual{
		address:     address,
		ordersRepo:  ordersRepo,
		balanceRepo: balanceRepo,
	}
}

func (a *Accrual) Sync() {
	for {
		ctx := context.Background()
		time.Sleep(5 * time.Second)
		
		newOrders, err := a.ordersRepo.GetNew(ctx)
		if err != nil {
			slog.Error("accrual sync", slog.String("method", "getNewOrders"), slog.String("error", err.Error()))
			continue
		}

		slog.Debug("accrual sync", slog.Any("new_orders", newOrders))

		for _, order := range newOrders {
			accrualData, err := withRetry(a.GetOrderStatus)(order.OrderID)
			if err != nil {
				slog.Error("accrual sync", slog.String("method", "getOrderStatus"), slog.String("error", err.Error()))
				continue
			}

			slog.Debug("accrual sync", slog.Any("accrual_data", *accrualData))

			if order.Status != accrualData.Status {
				userBalanceData, err := a.balanceRepo.GetByLogin(ctx, order.UserLogin)
				if err != nil {
					slog.Error("accrual sync", slog.String("method", "getBalance"), slog.String("error", err.Error()))
					if err != sql.ErrNoRows {
						continue
					}

					userBalanceData = &models.BalanceData{
						UserLogin: order.UserLogin,
					}
				}

				order.Status = accrualData.Status
				order.Accrual = accrualData.Accrual
				userBalanceData.Current += accrualData.Accrual

				if err = a.ordersRepo.Update(ctx, order); err != nil {
					slog.Error("accrual sync", slog.String("method", "ordersRepo.Update"), slog.String("error", err.Error()))
					continue
				}

				if err = a.balanceRepo.Update(ctx, userBalanceData); err != nil {
					slog.Error("accrual sync", slog.String("method", "balanceRepo.Update"), slog.String("error", err.Error()))
					continue
				}
			}
			slog.Debug("accrual sync. no new orders")
		}
	}
}

func (a *Accrual) GetOrderStatus(orderID string) (*models.AccrualOrders, error) {
	response, err := resty.New().R().Get(a.address + accrualRoot + orderID)
	if err != nil {
		return nil, fmt.Errorf("failed make request: %s", err.Error())
	}

	accrual := models.AccrualOrders{}

	if err = json.Unmarshal(response.Body(), &accrual); err != nil {
		return nil, fmt.Errorf("failed unmarshal request body: %s", err.Error())
	}

	return &accrual, nil
}

func withRetry(f func(string) (*models.AccrualOrders, error)) func(string) (*models.AccrualOrders, error) {
	return func(s string) (*models.AccrualOrders, error) {
		for range 5 {
			if res, err := f(s); err == nil {
				return res, nil
			}
			slog.Warn("failed get accrual order status, will retry in 5 seconds")
			time.Sleep(5 * time.Second)
		}
		return nil, fmt.Errorf("timed out")
	}
}
