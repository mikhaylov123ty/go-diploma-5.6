package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

func (p *Postgres) CreateOrder(order *models.OrderData) error {

	return nil
}

func (p *Postgres) GetOrders(userID string) ([]*models.OrderData, error) {

	return nil, nil
}

func (p *Postgres) GetOrderByID(orderID string) (*models.OrderData, error) {

	return nil, nil
}

func (p *Postgres) Update(data *models.OrderData) error {

	return nil
}

func (p *Postgres) Close() error {
	p.DB.Close()
	return nil
}
