package withdrawals

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Storage interface {
	Update(*models.WithdrawData) error
	Get(string) ([]*models.WithdrawData, error)
}
