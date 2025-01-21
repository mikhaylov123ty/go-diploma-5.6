package api

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
)

type RegisterHandler struct {
	userSaver           registerUserSaver
	transactionsHandler transactionsHandler
}

type registerUserSaver interface {
	Save(context.Context, string, string) error
}

func NewRegisterHandler(userSaver registerUserSaver, transactionsHandler transactionsHandler) *RegisterHandler {
	return &RegisterHandler{
		userSaver,
		transactionsHandler,
	}
}

func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)
	pass := r.Context().Value(utils.ContextKey("pass")).(string)

	//TODO fix error and salt pass

	slog.DebugContext(r.Context(), "register handler", slog.String("login", login), slog.String("pass", pass))
	if err := h.transactionsHandler.Begin(); err != nil {
		slog.Error("register handler", slog.String("method", "begin_transaction"), slog.String("errpr", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.userSaver.Save(r.Context(), login, pass); err != nil {
		_ = h.transactionsHandler.Rollback()
		slog.ErrorContext(r.Context(), "register handler", slog.String("method", "saveUser"), slog.String("error", err.Error()))
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.transactionsHandler.Commit(); err != nil {
		slog.Error("register handler", slog.String("method", "begin_transaction"), slog.String("errpr", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", r.Context().Value(utils.ContextKey("token")).(string))
	w.WriteHeader(http.StatusOK)
}
