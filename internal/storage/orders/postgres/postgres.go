package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
)

type Postgres struct {
	db *sql.DB
}

func Init(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) CreateOrder(order *models.OrderData) error {

	return nil
}

func (p *Postgres) GetOrders(userID string) ([]*models.OrderData, error) {

	return nil, nil
}

func (p *Postgres) GetNewOrders() ([]*models.OrderData, error) {

	return nil, nil
}

func (p *Postgres) GetOrderByID(orderID string) (*models.OrderData, error) {

	return nil, nil
}

func (p *Postgres) Update(data *models.OrderData) error {

	return nil
}
