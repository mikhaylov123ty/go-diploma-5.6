package postgres

import (
	"database/sql"
	"fmt"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Postgres struct {
	DB *sql.DB
}

func Init(dbURI string) (*Postgres, error) {
	conn, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed connect to db: %w", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed ping db: %w", err)
	}

	return &Postgres{conn}, nil
}

func (p *Postgres) Create(withdraw *models.WithdrawData) error {

	return nil
}

func (p *Postgres) Get(userlogin string) ([]*models.WithdrawData, error) {

	return nil, nil
}

func (p *Postgres) Close() error {
	p.DB.Close()
	return nil
}
