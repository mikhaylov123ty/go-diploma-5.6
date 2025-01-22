package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type Postgres struct {
	db          *sql.DB
	transaction *sql.Tx
}

func Init(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) Begin() error {
	var err error

	p.transaction, err = p.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	slog.Debug("starting transaction")

	return nil
}

func (p *Postgres) Commit() error {
	if err := p.transaction.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	slog.Debug("commiting transaction")

	return nil
}

func (p *Postgres) Rollback() error {
	if err := p.transaction.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	slog.Debug("rolling back transaction")

	return nil
}
