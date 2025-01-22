package users

import (
	"context"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Storage interface {
	GetByLogin(context.Context, string) (*models.UserData, error)
	Save(context.Context, string, string) error
}
