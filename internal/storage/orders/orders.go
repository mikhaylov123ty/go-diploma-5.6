package orders

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Storage interface {
	CreateOrder(*models.OrderData) error
	GetOrders(string) ([]*models.OrderData, error)
	GetNewOrders() ([]*models.OrderData, error)
	GetOrderByID(string) (*models.OrderData, error)
	Update(*models.OrderData) error
}
