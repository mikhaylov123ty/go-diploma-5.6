package memory

import (
	"context"
	"sync"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Memory struct {
	mu sync.RWMutex
	DB map[string]*models.WithdrawData
}

func Init() *Memory {
	return &Memory{
		DB: make(map[string]*models.WithdrawData),
	}
}

func (m *Memory) Update(ctx context.Context, withdraw *models.WithdrawData) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DB[withdraw.Order] = withdraw

	return nil
}

func (m *Memory) GetByLogin(ctx context.Context, userlogin string) ([]*models.WithdrawData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var data []*models.WithdrawData

	for _, withdrawal := range m.DB {
		if withdrawal.UserLogin == userlogin {
			data = append(data, withdrawal)
		}
	}

	return data, nil
}
