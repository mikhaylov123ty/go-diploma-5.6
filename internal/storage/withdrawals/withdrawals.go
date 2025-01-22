package withdrawals

import (
	"context"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Storage interface {
	Update(context.Context, *models.WithdrawData) error
	GetByLogin(context.Context, string) ([]*models.WithdrawData, error)
}
