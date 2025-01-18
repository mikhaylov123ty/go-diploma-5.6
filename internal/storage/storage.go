package storage

import (
	"database/sql"
	"fmt"
	
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/balance"
	balanceMemory "github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/balance/memory"
	balancePostgres "github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/balance/postgres"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/orders"
	ordersMemory "github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/orders/memory"
	ordersPostgres "github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/orders/postgres"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/users"
	usersMemory "github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/users/memory"
	usersPostgres "github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/users/postgres"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/withdrawals"
	withdrawalsMemory "github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/withdrawals/memory"
	withdrawalsPostgres "github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/withdrawals/postgres"

	_ "github.com/lib/pq"
)

type Storage struct {
	Conn           *sql.DB
	UsersRepo      users.Storage
	OrdersRepo     orders.Storage
	BalanceRepo    balance.Storage
	WithdrawalRepo withdrawals.Storage
}

func New(dbURI string) (*Storage, error) {
	if dbURI != "" {
		conn, err := sql.Open("postgres", dbURI)
		if err != nil {
			return nil, fmt.Errorf("failed connect to db: %w", err)
		}

		if err = conn.Ping(); err != nil {
			return nil, fmt.Errorf("failed ping db: %w", err)
		}

		return &Storage{
			UsersRepo:      usersPostgres.Init(conn),
			OrdersRepo:     ordersPostgres.Init(conn),
			BalanceRepo:    balancePostgres.Init(conn),
			WithdrawalRepo: withdrawalsPostgres.Init(conn),
		}, nil
	}

	return &Storage{
		UsersRepo:      usersMemory.Init(),
		OrdersRepo:     ordersMemory.Init(),
		BalanceRepo:    balanceMemory.Init(),
		WithdrawalRepo: withdrawalsMemory.Init(),
	}, nil
}
