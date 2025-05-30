package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Memory struct {
	mu sync.RWMutex
	DB map[string]*models.BalanceData
}

func Init() *Memory {
	return &Memory{
		DB: make(map[string]*models.BalanceData),
	}
}

func (m *Memory) GetByLogin(ctx context.Context, login string) (*models.BalanceData, error) {
	u, ok := m.DB[login]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil

}

func (m *Memory) Update(ctx context.Context, data *models.BalanceData) error {
	m.DB[data.UserLogin] = data
	return nil
}
