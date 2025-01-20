package accrual

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
		time.Sleep(5 * time.Second)
		newOrders, err := a.ordersRepo.GetNewOrders()
		if err != nil {
			log.Printf("Error syncing accrual orders: %v", err)
			continue
		}

		log.Printf("Syncing accrual orders, New Orders: %v", newOrders)
		for _, order := range newOrders {
			accrualData, err := a.GetOrderStatus(order.OrderID)
			if err != nil {
				log.Printf("failed get accrural order status: %v", err)
				continue
			}
			slog.Debug("Syncing accrual orders", slog.Any("Accrual Data", *accrualData))
			if order.Status != accrualData.Status {
				userBalanceData, err := a.balanceRepo.GetBalance(order.UserLogin)
				if err != nil {
					if err != sql.ErrNoRows {
						log.Printf("failed get user balance data: %v", err)
						continue
					}
					userBalanceData = &models.BalanceData{
						UserLogin: order.UserLogin,
						Current:   0,
						Withdrawn: 0,
					}
				}

				order.Status = accrualData.Status
				order.Accrual = accrualData.Accrual
				userBalanceData.Current += accrualData.Accrual

				if err = a.ordersRepo.Update(order); err != nil {
					log.Printf("Error saving accrual orders: %v", err)
					continue
				}

				if err = a.balanceRepo.Update(userBalanceData); err != nil {
					log.Printf("Error saving accrual orders: %v", err)
					continue
				}
			}
			log.Println("No new orders")
		}
	}
}

func (a *Accrual) GetOrderStatus(orderID string) (*models.AccrualOrders, error) {
	log.Printf("Getting order status for order: %s", a.address+accrualRoot+orderID)
	response, err := resty.New().R().Get(a.address + accrualRoot + orderID)
	if err != nil {
		log.Printf("Error syncing accrual orders: %v", err)
		return nil, err
	}
	fmt.Println("RESPONSE BODY", response.Body())
	accrual := models.AccrualOrders{}

	if err = json.Unmarshal(response.Body(), &accrual); err != nil {
		log.Printf("Error unmarshaling accrual orders: %v", err)
	}
	fmt.Println(accrual)
	return &accrual, nil
}
