package api

import (
	"encoding/json"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"log"
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
	GetBalance(string) (*models.BalanceData, error)
}

func NewGetBalanceHandler(balanceProvider balanceGetUserProvider) *BalanceHandler {
	return &BalanceHandler{
		balanceProvider,
	}
}

func (h *BalanceHandler) Handle(w http.ResponseWriter, r *http.Request) {
	balance, err := h.balanceProvider.GetBalance(r.Context().Value("login").(string))
	if err != nil {
		if err.Error() != "user not found" {
			log.Printf("error getting user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var resData BalanceResponse
	resData.Current = balance.Current
	resData.Withdrawn = balance.Withdrawn

	res, err := json.Marshal(resData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
