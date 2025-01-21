package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"

	"github.com/Masterminds/squirrel"
)

type Postgres struct {
	db *sql.DB
}

func Init(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) Update(ctx context.Context, data *models.WithdrawData) error {
	query, args, err := squirrel.Insert("withdrawals").
		Columns("order_id", "user_login", "sum", "processed_at").
		Values(data.Order, data.UserLogin, data.Sum, data.ProcessedAt).
		Suffix("ON CONFLICT (order_id) DO UPDATE").
		Suffix("SET sum = $3, processed_at = $4").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	slog.DebugContext(ctx, "Update Withdrawals", slog.String("query", query), slog.Any("args", args))

	res, err := p.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	if resAff, _ := res.RowsAffected(); resAff == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *Postgres) GetByLogin(ctx context.Context, login string) ([]*models.WithdrawData, error) {
	query, args, err := squirrel.Select("*").
		From("withdrawals").
		Where(squirrel.Eq{"user_login": login}).
		OrderBy("processed_at DESC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	slog.DebugContext(ctx, "Get New Orders", slog.String("query", query), slog.Any("args", args))

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	var res []*models.WithdrawData
	for rows.Next() {
		order := &models.WithdrawData{}
		if err = rows.Scan(&order.Order, &order.UserLogin, &order.Sum, &order.ProcessedAt); err != nil {
			return nil, err
		}
		res = append(res, order)
	}

	return res, nil
}
