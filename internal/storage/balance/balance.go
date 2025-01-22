package balance

import (
	"context"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Storage interface {
	GetByLogin(context.Context, string) (*models.BalanceData, error)
	Update(context.Context, *models.BalanceData) error
}
