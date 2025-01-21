package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
)

type WithdrawalsHandler struct {
	withdrawalsProvider withdrawalsProvider
	transactionsHandler transactionsHandler
}

type withdrawalsProvider interface {
	GetByLogin(context.Context, string) ([]*models.WithdrawData, error)
}

func NewWithdrawalsHandler(withdrawalsProvider withdrawalsProvider, transactionsHandler transactionsHandler) *WithdrawalsHandler {
	return &WithdrawalsHandler{
		withdrawalsProvider,
		transactionsHandler,
	}
}

func (h *WithdrawalsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)
	if login == "" {
		slog.ErrorContext(r.Context(), "get orders handler. empty login")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.transactionsHandler.Begin(); err != nil {
		slog.ErrorContext(r.Context(), "withdrawals handler", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resData, err := h.withdrawalsProvider.GetByLogin(r.Context(), login)
	if err != nil {
		_ = h.transactionsHandler.Rollback()
		slog.ErrorContext(r.Context(), "withdrawals handler",
			slog.String("method", "withdrawalsProvider.Get"),
			slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.transactionsHandler.Commit(); err != nil {
		slog.ErrorContext(r.Context(), "withdrawals handler", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(resData)
	if err != nil {
		slog.ErrorContext(r.Context(), "withdrawals handler",
			slog.String("method", "marshal response"),
			slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
