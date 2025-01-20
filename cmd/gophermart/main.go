package main

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/logger"
	"log"
	"log/slog"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/config"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server/accrual"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage"
)

// TODO logger, retrier, workers pool transactions and graceful shutdown

func main() {
	//init config
	cfg, err := config.Init()
	if err != nil {
		log.Fatal("failed init config: ", err)
	}

	logger.Init(cfg.LogLevel)

	slog.Info("test")

	//init storage
	storages, err := storage.New(cfg.DBURI)
	if err != nil {
		log.Fatal("failed init storages: ", err)
	}
	defer storages.Conn.Close()

	//init accural
	accrualInstance := accrual.NewAccrual(
		cfg.AccuralSystemAddress,
		storages.OrdersRepo,
		storages.BalanceRepo,
	)

	//init server
	serverInstance := server.New(
		cfg.Address,
		storages.UsersRepo,
		storages.OrdersRepo,
		storages.BalanceRepo,
		storages.WithdrawalRepo,
		cfg.Secret,
	)

	//start processing accrual orders
	go accrualInstance.Sync()

	//run server
	if err = serverInstance.Start(); err != nil {
		log.Fatal("failed start server: ", err)
	}
}
