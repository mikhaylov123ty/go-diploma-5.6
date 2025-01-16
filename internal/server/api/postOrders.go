package api

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"io"
	"log"
	"net/http"
	"time"
)

type OrderPostHandler struct {
	orderSaver   ordersPostSaver
	userProvider ordersPostUserProvider
}

type ordersPostSaver interface {
	CreateOrder(*models.OrderData) error
	GetOrderByID(string) (*models.OrderData, error)
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR READ BODY", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	login := r.Context().Value("login").(string)

	user, err := h.userProvider.GetUser(login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orderID := string(body)

	orderCheck, err := h.orderSaver.GetOrderByID(orderID)
	if err != nil && err.Error() != "order not found" {
		log.Println("error get order by id", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if orderCheck != nil {
		if orderCheck.UserLogin != user.Login {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	order := &models.OrderData{
		OrderID:    orderID,
		UserLogin:  user.Login,
		Status:     "NEW",
		UploadedAt: time.Now(),
	}

	if err = h.orderSaver.CreateOrder(order); err != nil {
		log.Println("ERROR CREATE ORDER", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("CREATED ORDER", orderID)

	w.WriteHeader(http.StatusAccepted)
}
