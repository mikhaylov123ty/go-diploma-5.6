package api

import (
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server/utils"
)

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"pass"`
}

type AuthHandler struct {
	userProvider loginUserProvider
}

type loginUserProvider interface {
	GetByLogin(context.Context, string) (*models.UserData, error)
}

func NewAuthHandler(userProvider loginUserProvider) *AuthHandler {
	return &AuthHandler{
		userProvider,
	}
}

func (h *AuthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)
	pass := r.Context().Value(utils.ContextKey("pass")).(string)

	user, err := h.userProvider.GetByLogin(r.Context(), login)
	if err != nil {
		slog.ErrorContext(r.Context(), "auth handler", slog.String("method", "getUser"), slog.String("error", err.Error()))

		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.DebugContext(r.Context(), "auth handler", slog.String("method", "getUser"), slog.Any("user", user))

	if err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(pass)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Authorization", r.Context().Value(utils.ContextKey("token")).(string))
	w.WriteHeader(http.StatusOK)
}
