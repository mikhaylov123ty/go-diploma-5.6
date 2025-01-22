package api

import (
	"context"
	"database/sql"

	"log/slog"
	"net/http"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server/utils"

	"golang.org/x/crypto/bcrypt"
)

type RegisterHandler struct {
	userSaver registerUserSaver
}

type registerUserSaver interface {
	Save(context.Context, string, string) error
}

func NewRegisterHandler(userSaver registerUserSaver) *RegisterHandler {
	return &RegisterHandler{
		userSaver,
	}
}

func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)
	pass := r.Context().Value(utils.ContextKey("pass")).(string)

	if login == "" {
		slog.ErrorContext(r.Context(), "register handler. empty login")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		slog.ErrorContext(r.Context(), "register handler. failed to hash password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.userSaver.Save(r.Context(), login, string(hashedPass)); err != nil {
		slog.ErrorContext(r.Context(), "register handler", slog.String("method", "saveUser"), slog.String("error", err.Error()))
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Authorization", r.Context().Value(utils.ContextKey("token")).(string))
	w.WriteHeader(http.StatusOK)
}
