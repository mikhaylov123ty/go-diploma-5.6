package api

import (
	"context"
	"database/sql"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
	"log/slog"
	"net/http"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"pass"`
}

type AuthHandler struct {
	userProvider        loginUserProvider
	transactionsHandler transactionsHandler
}

type loginUserProvider interface {
	GetByLogin(context.Context, string) (*models.UserData, error)
}

func NewAuthHandler(userProvider loginUserProvider, transactionsHandler transactionsHandler) *AuthHandler {
	return &AuthHandler{
		userProvider,
		transactionsHandler,
	}
}

func (h *AuthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)
	pass := r.Context().Value(utils.ContextKey("pass")).(string)

	if err := h.transactionsHandler.Begin(); err != nil {
		slog.ErrorContext(r.Context(), "auth handler", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.userProvider.GetByLogin(r.Context(), login)
	if err != nil {
		_ = h.transactionsHandler.Rollback()
		slog.ErrorContext(r.Context(), "auth handler", slog.String("method", "getUser"), slog.String("error", err.Error()))

		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.transactionsHandler.Commit(); err != nil {
		slog.ErrorContext(r.Context(), "auth handler", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.DebugContext(r.Context(), "auth handler", slog.String("method", "getUser"), slog.Any("user", user))
	if user.Pass != "" && user.Pass != pass {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Authorization", r.Context().Value(utils.ContextKey("token")).(string))
	w.WriteHeader(http.StatusOK)
}
