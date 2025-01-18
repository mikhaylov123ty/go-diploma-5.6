package main

import (
	"log"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/config"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server/accrual"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage"
)

func main() {
	//init config
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

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

	//start processing accrual orders
	go accrualInstance.Sync()

	//init server
	serverInstance := server.New(
		cfg.Address,
		storages.UsersRepo,
		storages.OrdersRepo,
		storages.BalanceRepo,
		storages.WithdrawalRepo,
		cfg.Secret,
	)

	//run server
	if err = serverInstance.Start(); err != nil {
		log.Fatal(err)
	}
}
