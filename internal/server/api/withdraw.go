package api

import (
	"database/sql"
	"encoding/json"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
	"io"
	"log"
	"net/http"
	"time"
)

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
type WithdrawHandler struct {
	balanceProvider  balanceProvider
	orderProvider    orderProvider
	withdrawProvider withdrawProvider
}

type balanceProvider interface {
	GetBalance(string) (*models.BalanceData, error)
	Update(*models.BalanceData) error
}

type orderProvider interface {
	GetOrderByID(string) (*models.OrderData, error)
	Update(*models.OrderData) error
}

type withdrawProvider interface {
	Update(*models.WithdrawData) error
}

func NewWithdrawHandler(balanceProvider balanceProvider, orderProvider orderProvider, withdrawProvider withdrawProvider) *WithdrawHandler {
	return &WithdrawHandler{
		balanceProvider:  balanceProvider,
		orderProvider:    orderProvider,
		withdrawProvider: withdrawProvider,
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

	userLogin := r.Context().Value(utils.ContextKey("login")).(string)
	log.Printf("WITHDRAW User login: %s", userLogin)

	order, err := h.orderProvider.GetOrderByID(req.Order)
	log.Printf("WITHDRAW Order ID: %s", req.Order)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("error getting order: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if order == nil {
		balance, err := h.balanceProvider.GetBalance(userLogin)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Printf("error getting user: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("WITHDRAW Balance: %v", balance)

		if balance.Current < req.Sum {
			log.Printf("balance is too low: %v", balance.Current)
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		balance.Current -= req.Sum
		balance.Withdrawn += req.Sum

		if err = h.balanceProvider.Update(balance); err != nil {
			log.Printf("error updating balance: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		withdrawData := &models.WithdrawData{
			UserLogin:   userLogin,
			Order:       req.Order,
			Sum:         req.Sum,
			ProcessedAt: time.Now(),
		}

		if err = h.withdrawProvider.Update(withdrawData); err != nil {
			log.Printf("error updating withdraw: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	log.Println("WITHDRAW order not found: ", err)
	w.WriteHeader(http.StatusUnprocessableEntity)
}
