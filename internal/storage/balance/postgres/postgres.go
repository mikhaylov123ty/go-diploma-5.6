package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log/slog"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func Init(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) GetByLogin(ctx context.Context, login string) (*models.BalanceData, error) {
	query, args, err := squirrel.Select("*").
		From("balances").
		Where(squirrel.Eq{"user_login": login}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	slog.DebugContext(ctx, "Get Balance", slog.String("query", query), slog.Any("args", args))

	row := p.db.QueryRow(query, args...)
	if row.Err() != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	var res models.BalanceData
	if err = row.Scan(&res.UserLogin, &res.Current, &res.Withdrawn); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("scanning row: %w", err)
	}

	return &res, nil
}

func (p *Postgres) Update(ctx context.Context, data *models.BalanceData) error {
	query, args, err := squirrel.Insert("balances").
		Columns("user_login", "current", "withdrawn").
		Values(data.UserLogin, data.Current, data.Withdrawn).
		Suffix("ON CONFLICT (user_login) DO UPDATE").
		Suffix("SET current = $2, withdrawn = $3").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	slog.DebugContext(ctx, "Update Balance", slog.String("query", query), slog.Any("args", args))

	res, err := p.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	if resAff, _ := res.RowsAffected(); resAff == 0 {
		return sql.ErrNoRows
	}

	return nil
}
