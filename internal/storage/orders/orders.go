package orders

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/orders/memory"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/orders/postgres"
)

type Storage interface {
	SaveOrder(string, string) error
	GetOrders(string) ([]*models.OrderData, error)
	Close() error
}

func New(dbURI string) (Storage, error) {
	if dbURI != "" {
		psgConn, err := postgres.Init(dbURI)
		if err != nil {
			return nil, err
		}

		return psgConn, nil
	}

	memoryConn, err := memory.Init()
	if err != nil {
		return nil, err
	}

	return memoryConn, nil
}
