package api

import (
	"encoding/json"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
	"log"
	"net/http"
)

type WithdrawalsHandler struct {
	withdrawalsProvider withdrawalsProvider
}

type withdrawalsProvider interface {
	Get(userlogin string) ([]*models.WithdrawData, error)
}

func NewWithdrawalsHandler(withdrawalsProvider withdrawalsProvider) *WithdrawalsHandler {
	return &WithdrawalsHandler{
		withdrawalsProvider,
	}
}

func (h *WithdrawalsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(utils.ContextKey("login")).(string)

	resData, err := h.withdrawalsProvider.Get(login)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(resData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
