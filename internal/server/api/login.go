package api

import (
	"database/sql"
	"fmt"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
	"log"
	"net/http"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"pass"`
}

type AuthHandler struct {
	provider userProvider
}

type userProvider interface {
	GetUser(string) (*models.UserData, error)
}

func NewAuthHandler(provider userProvider) *AuthHandler {
	return &AuthHandler{
		provider,
	}
}

func (h *AuthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)
	pass := r.Context().Value(utils.ContextKey("pass")).(string)
	user, err := h.provider.GetUser(login)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(user)

	if user.Pass != "" && user.Pass != pass {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Add("Authorization", r.Context().Value(utils.ContextKey("token")).(string))
	w.WriteHeader(http.StatusOK)
}
