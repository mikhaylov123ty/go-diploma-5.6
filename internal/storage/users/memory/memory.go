package memory

import (
	"context"
	"fmt"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"sync"
)

type Memory struct {
	mu sync.RWMutex
	DB map[string]*models.UserData
}

func Init() *Memory {
	return &Memory{
		DB: make(map[string]*models.UserData),
	}
}

func (m *Memory) Save(ctx context.Context, login string, pass string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.DB[login]; ok {
		return fmt.Errorf("user already exists")
	}

	m.DB[login] = &models.UserData{
		Login: login,
		Pass:  pass,
	}

	return nil
}

func (m *Memory) GetByLogin(ctx context.Context, login string) (*models.UserData, error) {
	u, ok := m.DB[login]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}
