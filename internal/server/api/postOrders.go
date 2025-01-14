package api

import (
	"encoding/json"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"io"
	"net/http"
)

type OrderPostHandler struct {
	orderSaver   ordersPostSaver
	userProvider ordersPostUserProvider
}

type ordersPostSaver interface {
	SaveOrder(string, string) error
}

type ordersPostUserProvider interface {
	GetUser(string) (*models.UserData, error)
}

func NewPostOrdersHandler(orderSaver ordersPostSaver, userProvider ordersPostUserProvider) *OrderPostHandler {
	return &OrderPostHandler{
		orderSaver,
		userProvider,
	}
}

func (h *OrderPostHandler) Handle(w http.ResponseWriter, r *http.Request) {
	login, err := r.Cookie("login")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userProvider.GetUser(login.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var orderID string
	if err = json.Unmarshal(body, &orderID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.orderSaver.SaveOrder(user.Login, orderID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
