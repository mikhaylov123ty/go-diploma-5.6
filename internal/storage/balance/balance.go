package balance

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/balance/memory"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/balance/postgres"
)

type Storage interface {
	GetBalance(string) (*models.BalanceData, error)
	Update(*models.BalanceData) error
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
