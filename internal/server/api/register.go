package api

import (
	"database/sql"
	"log"
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
	login := (r.Context().Value("login")).(string)
	pass := r.Context().Value("pass").(string)
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

	w.Header().Add("Authorization", r.Context().Value("token").(string))
	w.WriteHeader(http.StatusOK)
}
