package api

import (
	"encoding/json"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
	"net/http"
)

type OrderGetHandler struct {
	ordersProvider ordersGetProvider
	userProvider   ordersGetUserProvider
}

type ordersGetProvider interface {
	GetOrders(string) ([]*models.OrderData, error)
}

type ordersGetUserProvider interface {
	GetUser(string) (*models.UserData, error)
}

func NewGetOrdersHandler(ordersProvider ordersGetProvider, userProvider ordersGetUserProvider) *OrderGetHandler {
	return &OrderGetHandler{
		ordersProvider,
		userProvider,
	}
}

func (h *OrderGetHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)
	user, err := h.userProvider.GetUser(login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orders, err := h.ordersProvider.GetOrders(user.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
