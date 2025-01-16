package memory

import (
	"fmt"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"sync"
)

type Memory struct {
	mu sync.RWMutex
	DB map[string]*models.OrderData
}

func Init() (*Memory, error) {
	return &Memory{
		DB: make(map[string]*models.OrderData),
	}, nil
}

func (m *Memory) CreateOrder(order *models.OrderData) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DB[order.OrderID] = order

	return nil
}

func (m *Memory) GetOrders(userlogin string) ([]*models.OrderData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var orders []*models.OrderData

	for _, order := range m.DB {
		if order.UserLogin == userlogin {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (m *Memory) GetOrderByID(orderID string) (*models.OrderData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.DB[orderID]; !ok {
		return nil, fmt.Errorf("order not found")
	}

	return m.DB[orderID], nil
}

func (m *Memory) Update(data *models.OrderData) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DB[data.OrderID] = data

	return nil
}

func (m *Memory) Close() error {
	return nil
}
