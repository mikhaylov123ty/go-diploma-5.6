package api

import (
	"context"
	"database/sql"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type OrderPostHandler struct {
	orderSaver   ordersPostSaver
	userProvider ordersPostUserProvider
}

type ordersPostSaver interface {
	Create(context.Context, *models.OrderData) error
	GetByID(context.Context, string) (*models.OrderData, error)
}

type ordersPostUserProvider interface {
	GetByLogin(context.Context, string) (*models.UserData, error)
}

func NewPostOrdersHandler(orderSaver ordersPostSaver, userProvider ordersPostUserProvider) *OrderPostHandler {
	return &OrderPostHandler{
		orderSaver,
		userProvider,
	}
}

func (h *OrderPostHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(r.Context(), "order post handler", slog.String("method", "read_body"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	login := r.Context().Value(utils.ContextKey("login")).(string)
	if login == "" {
		slog.ErrorContext(r.Context(), "order post handler. empty login")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userProvider.GetByLogin(r.Context(), login)
	if err != nil {
		slog.ErrorContext(r.Context(), "order post handler", slog.String("method", "getUser"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user == nil {
		slog.DebugContext(r.Context(), "order post handler", slog.String("method", "user_nill_check"), slog.Bool("is nill", false))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orderID := string(body)
	if !checkLuhn(orderID) {
		slog.DebugContext(r.Context(), "order post handler", slog.String("method", "Luhn_order_check"), slog.Bool("is ok", false))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	orderCheck, err := h.orderSaver.GetByID(r.Context(), orderID)
	if err != nil && err != sql.ErrNoRows {
		slog.ErrorContext(r.Context(), "order post handler", slog.String("method", "getOrderByID"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if orderCheck != nil {
		if orderCheck.UserLogin != user.Login {
			slog.DebugContext(r.Context(), "order post handler", slog.String("method", "user_check"), slog.Bool("is ok", false))
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	order := &models.OrderData{
		OrderID:    orderID,
		UserLogin:  user.Login,
		Status:     "NEW",
		UploadedAt: time.Now(),
	}

	if err = h.orderSaver.Create(r.Context(), order); err != nil {
		slog.ErrorContext(r.Context(), "order post handler", slog.String("method", "createOrder"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.DebugContext(r.Context(), "order post handler", slog.Any("created_order", order))

	w.WriteHeader(http.StatusAccepted)
}

func checkLuhn(purportedCC string) bool {
	var sum = 0
	var parity = len(purportedCC) % 2

	for i, v := range purportedCC {
		v -= 48
		if i%2 == parity {
			v *= 2
			if v > 9 {
				v -= 9
			}
		}
		sum += int(v)
	}
	return sum%10 == 0
}
