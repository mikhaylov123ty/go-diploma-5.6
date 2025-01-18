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

func (p *Postgres) GetBalance(login string) (*models.BalanceData, error) {

	return &models.BalanceData{}, nil
}

func (P *Postgres) Update(data *models.BalanceData) error {

	return nil
}

func (p *Postgres) Close() error {
	p.db.Close()
	return nil
}
