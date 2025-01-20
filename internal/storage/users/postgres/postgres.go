package postgres

import (
	"database/sql"
	"fmt"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"

	"github.com/Masterminds/squirrel"
)

type Postgres struct {
	db *sql.DB
}

func Init(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) SaveUser(login string, pass string) error {
	query, args, err := squirrel.Insert("users").Columns("login", "pass").
		Values(login, pass).Suffix("ON CONFLICT (login) DO NOTHING").
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

func (p *Postgres) GetUser(login string) (*models.UserData, error) {
	query, args, err := squirrel.Select("*").
		From("users").
		Where(squirrel.Eq{"login": login}).
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

	var res models.UserData
	if err = row.Scan(&res.Login, &res.Pass); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("scanning row: %w", err)
	}

	return &res, nil
}
