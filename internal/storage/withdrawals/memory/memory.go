package memory

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"sync"
)

type Memory struct {
	mu sync.RWMutex
	DB map[string]*models.WithdrawData
}

func Init() (*Memory, error) {
	return &Memory{
		DB: make(map[string]*models.WithdrawData),
	}, nil
}

func (m *Memory) Create(withdraw *models.WithdrawData) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DB[withdraw.Order] = withdraw

	return nil
}

func (m *Memory) Get(userlogin string) ([]*models.WithdrawData, error) {
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

func (m *Memory) Close() error {
	return nil
}
