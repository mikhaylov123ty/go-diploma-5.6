package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
)

type OrdersGetHandler struct {
	ordersProvider ordersGetProvider
	userProvider   ordersGetUserProvider
}

type ordersGetProvider interface {
	GetByLogin(context.Context, string) ([]*models.OrderData, error)
}

type ordersGetUserProvider interface {
	GetByLogin(context.Context, string) (*models.UserData, error)
}

func NewGetOrdersHandler(ordersProvider ordersGetProvider, userProvider ordersGetUserProvider) *OrdersGetHandler {
	return &OrdersGetHandler{
		ordersProvider,
		userProvider,
	}
}

func (h *OrdersGetHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)
	if login == "" {
		slog.ErrorContext(r.Context(), "get orders handler. empty login")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userProvider.GetByLogin(r.Context(), login)
	if err != nil {
		slog.ErrorContext(r.Context(), "get orders handler", slog.String("method", "getUser"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orders, err := h.ordersProvider.GetByLogin(r.Context(), user.Login)
	if err != nil {
		slog.ErrorContext(r.Context(), "get orders handler", slog.String("method", "getOrders"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(orders)
	if err != nil {
		slog.ErrorContext(r.Context(), "get orders handler", slog.String("method", "marshal_orders"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
