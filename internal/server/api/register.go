package api

import (
	"database/sql"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
	"log"
	"log/slog"
	"net/http"
)

type RegisterHandler struct {
	userRegister userRegister
}

type userRegister interface {
	SaveUser(string, string) error
}

func NewRegisterHandler(userRegister userRegister) *RegisterHandler {
	return &RegisterHandler{
		userRegister,
	}
}

func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)
	pass := r.Context().Value(utils.ContextKey("pass")).(string)
	slog.DebugContext(r.Context(), "context login and pass", slog.String("login", login), slog.String("pass", pass))
	if err := h.userRegister.SaveUser(login, pass); err != nil {
		log.Println(err)
		//TODO fix error and salt pass
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
