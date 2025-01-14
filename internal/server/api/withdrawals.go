package api

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"net/http"
)

type WithdrawalsHandler struct {
	userProvider ordersGetUserProvider
}

type withdrawalsUserProvider interface {
	GetUser(string) (*models.UserData, error)
	SaveUser(string, string) error
}

func NewWithdrawalsHandler(withdrawalsUserProvider withdrawalsUserProvider) *WithdrawalsHandler {
	return &WithdrawalsHandler{
		withdrawalsUserProvider,
	}
}

func (h *WithdrawalsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	//TODO
	// get user

	//do something

	//save user
}
