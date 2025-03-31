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

func (p *Postgres) Save(ctx context.Context, login string, pass string) error {
	query, args, err := squirrel.Insert("users").Columns("login", "pass").
		Values(login, pass).Suffix("ON CONFLICT (login) DO NOTHING").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	slog.DebugContext(ctx, "Save User", slog.String("query", query), slog.Any("args", args))

	res, err := p.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	if resAff, _ := res.RowsAffected(); resAff == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *Postgres) GetByLogin(ctx context.Context, login string) (*models.UserData, error) {
	query, args, err := squirrel.Select("*").
		From("users").
		Where(squirrel.Eq{"login": login}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "Get User By Login", slog.String("query", query), slog.Any("args", args))

	row := p.db.QueryRow(query, args...)
	if row.Err() != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	var res models.UserData
	if err = row.Scan(&res.Login, &res.Pass); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("scanning row: %w", err)
	}

	return &res, nil
}
