package balance

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Storage interface {
	GetBalance(string) (*models.BalanceData, error)
	Update(*models.BalanceData) error
}
