package postgres

import (
	"database/sql"
	"fmt"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func Init(dbURI string) (*Postgres, error) {
	conn, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed connect to db: %w", err)
	}

	return &Postgres{conn}, nil
}

func (p *Postgres) SaveUser(login string, pass string) error {

	return nil
}

func (p *Postgres) GetUser(login string) (*models.UserData, error) {

	return nil, nil
}

func (p *Postgres) Close() error {
	p.DB.Close()
	return nil
}
