package users

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/users/memory"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/users/postgres"
)

type Storage interface {
	GetUser(string) (*models.UserData, error)
	SaveUser(string, string) error
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
