package accrual

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/transactions"

	"log/slog"
	"net/http"
	"time"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/balance"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/orders"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"

	"github.com/go-resty/resty/v2"
)

const (
	accrualRoot  = "/api/orders/"
	syncInterval = time.Second * 5
	retryCount   = 3
	retryDelay   = time.Second * 2
	workersPool  = 4
)

type Accrual struct {
	address     string
	ordersRepo  orders.Storage
	balanceRepo balance.Storage
	transaction transactions.Handler
}

func NewAccrual(
	address string,
	ordersRepo orders.Storage,
	balanceRepo balance.Storage,
	transaction transactions.Handler) *Accrual {
	return &Accrual{
		address:     address,
		ordersRepo:  ordersRepo,
		balanceRepo: balanceRepo,
		transaction: transaction,
	}
}

func (a *Accrual) Sync() {
	jobs := make(chan *models.OrderData)
	res := make(chan *utils.WorkersResponse)

	for range workersPool {
		go a.SyncWorker(jobs, res)
	}

	for {
		slog.Debug("accrual sync start loop")
		ctx := context.Background()
		time.Sleep(syncInterval)

		if err := a.transaction.Begin(); err != nil {
			slog.Error("accrual transaction begin err", slog.String("error", err.Error()))
			continue
		}

		newOrders, err := a.ordersRepo.GetNew(ctx)
		if err != nil {
			_ = a.transaction.Rollback()
			slog.Error("accrual sync", slog.String("method", "getNewOrders"), slog.String("error", err.Error()))
			continue
		}

		if err := a.transaction.Commit(); err != nil {
			slog.Error("accrual transaction commit err", slog.String("error", err.Error()))
		}

		if len(newOrders) < 1 {
			slog.Info("accrual sync. no new orders in orders repo")
			continue
		}

		slog.Debug("accrual sync", slog.Any("new_orders", newOrders))

		go func(newOrders []*models.OrderData) {
			for _, order := range newOrders {
				jobs <- order
			}
		}(newOrders)

		// Чтение результатов из результирующего канала
		for range len(newOrders) {
			r := <-res
			if r.Err != nil {
				slog.Error("worker failed processing accrual", slog.String("error", r.Err.Error()))
				continue
			}
			slog.Info("Accrual processed", slog.Any("Worker ID", r.WorkerID))
		}
	}
}

func (a *Accrual) GetOrderStatus(orderID string) (*models.AccrualOrders, error) {
	slog.Debug("get accrual order status request", slog.String("address", a.address+accrualRoot+orderID))
	response, err := resty.New().R().Get(a.address + accrualRoot + orderID)
	if err != nil {
		return nil, fmt.Errorf("failed make request: %s", err.Error())
	}
	if response.StatusCode() == http.StatusNoContent {
		return nil, nil
	}

	accrual := models.AccrualOrders{}
	if err = json.Unmarshal(response.Body(), &accrual); err != nil {
		return nil, fmt.Errorf("failed unmarshal request body: %s", err.Error())
	}

	return &accrual, nil
}

func withRetry(f func(string) (*models.AccrualOrders, error)) func(string) (*models.AccrualOrders, error) {
	return func(s string) (*models.AccrualOrders, error) {
		for range retryCount {
			res, err := f(s)
			if err == nil {
				return res, nil
			}
			slog.Warn("failed get accrual order status, will retry in 2 seconds", slog.String("error", err.Error()))
			time.Sleep(retryDelay)
		}
		return nil, fmt.Errorf("timed out")
	}
}

func (a *Accrual) SyncWorker(jobs <-chan *models.OrderData, results chan<- *utils.WorkersResponse) {
	for order := range jobs {
		ctx := context.Background()
		res := &utils.WorkersResponse{}

		accrualData, err := withRetry(a.GetOrderStatus)(order.OrderID)
		if err != nil {
			slog.Error("accrual sync", slog.String("method", "getOrderStatus"), slog.String("error", err.Error()))
			res.Err = err
			results <- res
			continue
		}
		if accrualData == nil {
			slog.Warn("accrual sync. no order found", slog.Any("order", order))
			results <- res
			continue
		}

		slog.Debug("accrual sync", slog.Any("accrual_data", *accrualData))

		if order.Status != accrualData.Status {
			if err = a.transaction.Begin(); err != nil {
				slog.Error("accrual transaction begin err", slog.String("error", err.Error()))
				res.Err = err
				results <- res

				continue
			}

			userBalanceData, err := a.balanceRepo.GetByLogin(ctx, order.UserLogin)
			if err != nil {
				slog.Error("accrual sync", slog.String("method", "getBalance"), slog.String("error", err.Error()))
				if err != sql.ErrNoRows {
					_ = a.transaction.Rollback()
					res.Err = err
					results <- res

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
				_ = a.transaction.Rollback()
				slog.Error("accrual sync", slog.String("method", "ordersRepo.Update"), slog.String("error", err.Error()))
				res.Err = err
				results <- res

				continue
			}

			if err = a.balanceRepo.Update(ctx, userBalanceData); err != nil {
				_ = a.transaction.Rollback()
				slog.Error("accrual sync", slog.String("method", "balanceRepo.Update"), slog.String("error", err.Error()))
				res.Err = err
				results <- res

				continue
			}

			if err = a.transaction.Commit(); err != nil {
				slog.Error("accural sync failed commit transaction", slog.String("error", err.Error()))
				res.Err = err
				results <- res

				continue
			}

			slog.Info("accrual sync", slog.Any("saved order", order))
			slog.Info("accrual sync", slog.Any("saved balance", userBalanceData))

			results <- res
		}
	}
}
