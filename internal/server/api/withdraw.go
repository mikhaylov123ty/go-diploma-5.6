package api

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"net/http"
)

type WithdrawHandler struct {
	userProvider ordersGetUserProvider
}

type withdrawGetUserProvider interface {
	GetUser(string) (*models.UserData, error)
}

func NewWithdrawHandler(withdrawGetUserProvider withdrawGetUserProvider) *WithdrawHandler {
	return &WithdrawHandler{
		withdrawGetUserProvider,
	}
}

func (h *WithdrawHandler) Handle(w http.ResponseWriter, r *http.Request) {
	//TODO
}
