package users

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Storage interface {
	GetUser(string) (*models.UserData, error)
	SaveUser(string, string) error
}
