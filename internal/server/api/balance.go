package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server/utils"
	"log/slog"
	"net/http"
)

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type BalanceHandler struct {
	balanceProvider balanceGetUserProvider
}

type balanceGetUserProvider interface {
	GetByLogin(context.Context, string) (*models.BalanceData, error)
}

func NewGetBalanceHandler(balanceProvider balanceGetUserProvider) *BalanceHandler {
	return &BalanceHandler{
		balanceProvider,
	}
}

func (h *BalanceHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)
	if login == "" {
		slog.ErrorContext(r.Context(), "balance handler. empty login")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	balance, err := h.balanceProvider.GetByLogin(r.Context(), login)
	if err != nil {
		slog.ErrorContext(r.Context(), "balance handler", slog.String("method", "get balance"), slog.String("error", err.Error()))
		if err != sql.ErrNoRows {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		return
	}

	var resData BalanceResponse
	resData.Current = balance.Current
	resData.Withdrawn = balance.Withdrawn

	res, err := json.Marshal(resData)
	if err != nil {
		slog.ErrorContext(r.Context(), "balance handler", slog.String("method", "marshal response data"), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
