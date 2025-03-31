package orders

import (
	"context"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Storage interface {
	Create(context.Context, *models.OrderData) error
	GetByLogin(context.Context, string) ([]*models.OrderData, error)
	GetNew(context.Context) ([]*models.OrderData, error)
	GetByID(context.Context, string) (*models.OrderData, error)
	Update(context.Context, *models.OrderData) error
}
