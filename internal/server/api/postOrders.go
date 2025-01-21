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
	orderSaver          ordersPostSaver
	userProvider        ordersPostUserProvider
	transactionsHandler transactionsHandler
}

type ordersPostSaver interface {
	Create(context.Context, *models.OrderData) error
	GetByID(context.Context, string) (*models.OrderData, error)
}

type ordersPostUserProvider interface {
	GetByLogin(context.Context, string) (*models.UserData, error)
}

func NewPostOrdersHandler(
	orderSaver ordersPostSaver,
	userProvider ordersPostUserProvider,
	transactionsHandler transactionsHandler) *OrderPostHandler {
	return &OrderPostHandler{
		orderSaver,
		userProvider,
		transactionsHandler,
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

	if err = h.transactionsHandler.Begin(); err != nil {
		slog.ErrorContext(r.Context(), "order post handler", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.userProvider.GetByLogin(r.Context(), login)
	if err != nil {
		_ = h.transactionsHandler.Rollback()
		slog.ErrorContext(r.Context(), "order post handler", slog.String("method", "getUser"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user == nil {
		_ = h.transactionsHandler.Rollback()
		slog.DebugContext(r.Context(), "order post handler", slog.String("method", "user_nill_check"), slog.Bool("is nill", false))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orderID := string(body)
	if !checkLuhn(orderID) {
		_ = h.transactionsHandler.Rollback()
		slog.DebugContext(r.Context(), "order post handler", slog.String("method", "Luhn_order_check"), slog.Bool("is ok", false))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	orderData, err := h.orderSaver.GetByID(r.Context(), orderID)
	if err != nil && err != sql.ErrNoRows {
		_ = h.transactionsHandler.Rollback()
		slog.ErrorContext(r.Context(), "order post handler", slog.String("method", "getOrderByID"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if orderData != nil {
		_ = h.transactionsHandler.Rollback()
		if orderData.UserLogin != user.Login {
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
		_ = h.transactionsHandler.Rollback()
		slog.ErrorContext(r.Context(), "order post handler", slog.String("method", "createOrder"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.transactionsHandler.Commit(); err != nil {
		slog.ErrorContext(r.Context(), "order post handler", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
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
