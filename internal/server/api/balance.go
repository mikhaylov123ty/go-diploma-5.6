package api

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"net/http"
)

type BalanceHandler struct {
	userProvider ordersGetUserProvider
}

type balanceGetUserProvider interface {
	GetUser(string) (*models.UserData, error)
}

func NewGetBalanceHandler(balanceGetUserProvider balanceGetUserProvider) *BalanceHandler {
	return &BalanceHandler{
		balanceGetUserProvider,
	}
}

func (h *BalanceHandler) Handle(w http.ResponseWriter, r *http.Request) {
	//TODO
}
