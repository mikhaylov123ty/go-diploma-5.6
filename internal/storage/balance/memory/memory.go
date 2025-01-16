package memory

import (
	"fmt"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"sync"
)

type Memory struct {
	mu sync.RWMutex
	DB map[string]*models.BalanceData
}

func Init() (*Memory, error) {
	return &Memory{
		DB: make(map[string]*models.BalanceData),
	}, nil
}

func (m *Memory) GetBalance(login string) (*models.BalanceData, error) {
	u, ok := m.DB[login]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil

}

func (m *Memory) Update(data *models.BalanceData) error {

	return nil
}

func (m *Memory) Close() error {
	return nil
}
