package api

import (
	"encoding/json"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"io"
	"log"
	"net/http"
)

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
type WithdrawHandler struct {
	balanceProvider balanceProvider
	orderProvider   orderProvider
}

type balanceProvider interface {
	GetBalance(string) (*models.BalanceData, error)
	Update(*models.BalanceData) error
}

type orderProvider interface {
	GetOrderByID(string) (*models.OrderData, error)
	Update(*models.OrderData) error
}

func NewWithdrawHandler(balanceProvider balanceProvider, orderProvider orderProvider) *WithdrawHandler {
	return &WithdrawHandler{
		balanceProvider: balanceProvider,
		orderProvider:   orderProvider,
	}
}

func (h *WithdrawHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR READ BODY", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req WithdrawRequest
	if err = json.Unmarshal(body, &req); err != nil {
		log.Println("ERROR READ BODY", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := h.orderProvider.GetOrderByID(req.Order)
	if err != nil {
		if err.Error() != "order not found" {
			log.Printf("error getting order: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	balance, err := h.balanceProvider.GetBalance(r.Context().Value("login").(string))
	if err != nil {
		if err.Error() != "user not found" {
			log.Printf("error getting user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if balance.Current < req.Sum {
		log.Printf("balance is too low: %v", balance.Current)
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	order.Accrual = req.Sum

	balance.Current -= req.Sum
	balance.Withdrawn += req.Sum

	if err = h.orderProvider.Update(order); err != nil {
		log.Printf("error updating order: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = h.balanceProvider.Update(balance); err != nil {
		log.Printf("error updating balance: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
