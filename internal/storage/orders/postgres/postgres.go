package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
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
	query, args, err := squirrel.Insert("orders").
		Columns("number", "user_login", "status", "accrual", "uploaded_at").
		Values(order.OrderID, order.UserLogin, order.Status, order.Accrual, order.UploadedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}
	fmt.Println("QUERY", query, "ARGS", args)

	res, err := p.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	if resAff, _ := res.RowsAffected(); resAff == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *Postgres) GetOrders(userID string) ([]*models.OrderData, error) {
	query, args, err := squirrel.Select("*").
		From("orders").
		Where(squirrel.Eq{"user_login": userID}).
		OrderBy("uploaded_at DESC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println("QUERY", query, "ARGS", args)

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	var res []*models.OrderData
	for rows.Next() {
		order := &models.OrderData{}
		if err = rows.Scan(&order.OrderID, &order.UserLogin, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		res = append(res, order)
	}

	return res, nil
}

func (p *Postgres) GetNewOrders() ([]*models.OrderData, error) {
	query, args, err := squirrel.Select("*").
		From("orders").
		Where(squirrel.Eq{"status": "NEW"}).
		OrderBy("uploaded_at DESC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println("QUERY", query, "ARGS", args)

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	var res []*models.OrderData
	for rows.Next() {
		order := &models.OrderData{}
		if err = rows.Scan(&order.OrderID, &order.UserLogin, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		res = append(res, order)
	}

	return res, nil
}

func (p *Postgres) GetOrderByID(orderID string) (*models.OrderData, error) {
	query, args, err := squirrel.Select("*").
		From("orders").
		Where(squirrel.Eq{"number": orderID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println("QUERY", query, "ARGS", args)

	row := p.db.QueryRow(query, args...)
	if row.Err() != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	var res models.OrderData
	if err = row.Scan(
		&res.OrderID,
		&res.UserLogin,
		&res.Status,
		&res.Accrual,
		&res.UploadedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("scanning row: %w", err)
	}

	return &res, nil

}

func (p *Postgres) Update(data *models.OrderData) error {
	query, args, err := squirrel.Insert("orders").
		Columns("number", "user_login", "status", "accrual", "uploaded_at").
		Values(data.OrderID, data.UserLogin, data.Status, data.Accrual, data.UploadedAt).
		Suffix("ON CONFLICT (number) DO UPDATE").
		Suffix("SET status = $3, accrual = $4").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	fmt.Println("QUERY", query, "ARGS", args)

	res, err := p.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	if resAff, _ := res.RowsAffected(); resAff == 0 {
		return sql.ErrNoRows
	}

	return nil
}
