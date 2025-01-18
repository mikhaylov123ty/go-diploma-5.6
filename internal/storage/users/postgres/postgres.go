package postgres

import (
	"database/sql"
	
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func Init(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) SaveUser(login string, pass string) error {

	return nil
}

func (p *Postgres) GetUser(login string) (*models.UserData, error) {

	return nil, nil
}
