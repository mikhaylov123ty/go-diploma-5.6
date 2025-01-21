package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
)

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawHandler struct {
	balanceProvider     balanceProvider
	orderProvider       orderProvider
	withdrawProvider    withdrawProvider
	transactionsHandler transactionsHandler
}

type balanceProvider interface {
	GetByLogin(context.Context, string) (*models.BalanceData, error)
	Update(context.Context, *models.BalanceData) error
}

type orderProvider interface {
	GetByID(context.Context, string) (*models.OrderData, error)
	Update(context.Context, *models.OrderData) error
}

type withdrawProvider interface {
	Update(context.Context, *models.WithdrawData) error
}

func NewWithdrawHandler(
	balanceProvider balanceProvider,
	orderProvider orderProvider,
	withdrawProvider withdrawProvider,
	transactionsHandler transactionsHandler) *WithdrawHandler {
	return &WithdrawHandler{
		balanceProvider:     balanceProvider,
		orderProvider:       orderProvider,
		withdrawProvider:    withdrawProvider,
		transactionsHandler: transactionsHandler,
	}
}

func (h *WithdrawHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(r.Context(), "withdraw handler",
			slog.String("method", "read_body"),
			slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req WithdrawRequest
	if err = json.Unmarshal(body, &req); err != nil {
		slog.ErrorContext(r.Context(), "withdraw handler",
			slog.String("method", "unmarshal_body"),
			slog.String("error", err.Error()))

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := r.Context().Value(utils.ContextKey("login")).(string)
	if login == "" {
		slog.ErrorContext(r.Context(), "order post handler. empty login")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.transactionsHandler.Begin(); err != nil {
		slog.ErrorContext(r.Context(), "withdraw handler", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	order, err := h.orderProvider.GetByID(r.Context(), req.Order)
	if err != nil && err != sql.ErrNoRows {
		_ = h.transactionsHandler.Rollback()
		slog.ErrorContext(r.Context(), "withdraw handler",
			slog.String("method", "getOrderByID"),
			slog.String("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.DebugContext(r.Context(), "withdraw handler", slog.Any("order", order))

	if order == nil {
		balance, err := h.balanceProvider.GetByLogin(r.Context(), login)
		if err != nil {
			_ = h.transactionsHandler.Rollback()
			slog.ErrorContext(r.Context(), "withdraw handler",
				slog.String("method", "getBalance"),
				slog.String("error", err.Error()))

			if err != sql.ErrNoRows {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNotFound)
			return
		}

		slog.DebugContext(r.Context(), "withdraw handler", slog.Any("balance", balance))

		if balance.Current < req.Sum {
			_ = h.transactionsHandler.Rollback()
			slog.DebugContext(r.Context(), "withdraw handler", slog.Any("balance is too low", balance.Current))
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		balance.Current -= req.Sum
		balance.Withdrawn += req.Sum

		if err = h.balanceProvider.Update(r.Context(), balance); err != nil {
			_ = h.transactionsHandler.Rollback()
			slog.ErrorContext(r.Context(), "withdraw handler",
				slog.String("method", "balanceProvider.Update"),
				slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		withdrawData := &models.WithdrawData{
			UserLogin:   login,
			Order:       req.Order,
			Sum:         req.Sum,
			ProcessedAt: time.Now(),
		}

		if err = h.withdrawProvider.Update(r.Context(), withdrawData); err != nil {
			_ = h.transactionsHandler.Rollback()
			slog.ErrorContext(r.Context(), "withdraw handler",
				slog.String("method", "withdrawProvider.Update"),
				slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = h.transactionsHandler.Commit(); err != nil {
			slog.ErrorContext(r.Context(), "withdraw handler", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	slog.DebugContext(r.Context(), "withdraw handler. order not found")

	w.WriteHeader(http.StatusUnprocessableEntity)
}
