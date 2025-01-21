package api

import (
	"context"
	"encoding/json"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
	"log/slog"
	"net/http"
)

type WithdrawalsHandler struct {
	withdrawalsProvider withdrawalsProvider
}

type withdrawalsProvider interface {
	GetByLogin(context.Context, string) ([]*models.WithdrawData, error)
}

func NewWithdrawalsHandler(withdrawalsProvider withdrawalsProvider) *WithdrawalsHandler {
	return &WithdrawalsHandler{
		withdrawalsProvider,
	}
}

func (h *WithdrawalsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)

	resData, err := h.withdrawalsProvider.GetByLogin(r.Context(), login)
	if err != nil {
		slog.ErrorContext(r.Context(), "withdrawals handler",
			slog.String("method", "withdrawalsProvider.Get"),
			slog.String("error", err.Error()))
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
