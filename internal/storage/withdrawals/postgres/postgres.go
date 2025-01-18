package postgres

import (
	"database/sql"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Postgres struct {
	db *sql.DB
}

func Init(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) Update(*models.WithdrawData) error {

	return nil
}

func (p *Postgres) Get(userlogin string) ([]*models.WithdrawData, error) {

	return nil, nil
}
