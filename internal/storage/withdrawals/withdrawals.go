package withdrawals

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/withdrawals/memory"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/withdrawals/postgres"
)

type Storage interface {
	Create(*models.WithdrawData) error
	Get(string) ([]*models.WithdrawData, error)
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
