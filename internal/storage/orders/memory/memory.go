package memory

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"sync"
)

type Memory struct {
	mu sync.RWMutex
	DB map[string]any
}

func Init() (*Memory, error) {
	return &Memory{
		DB: make(map[string]any),
	}, nil
}

func (m *Memory) SaveOrder(userID string, orderID string) error {

	return nil
}

func (m *Memory) GetOrders(userID string) ([]*models.OrderData, error) {

	return nil, nil
}

func (m *Memory) Close() error {
	return nil
}
