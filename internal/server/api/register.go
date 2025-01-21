package api

import (
	"context"
	"database/sql"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
	"log/slog"
	"net/http"
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

	//TODO fix error and salt pass

	slog.DebugContext(r.Context(), "register handler", slog.String("login", login), slog.String("pass", pass))
	if err := h.userSaver.Save(r.Context(), login, pass); err != nil {
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
