package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/config"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/logger"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server/accrual"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage"
)

func main() {
	//init config
	cfg, err := config.Init()
	if err != nil {
		log.Fatal("failed init config: ", err)
	}

	logger.Init(cfg.LogLevel)

	//init storage
	storages, err := storage.New(cfg.DBURI)
	if err != nil {
		log.Fatal("failed init storages: ", err)
	}

	//init accural
	accrualInstance := accrual.NewAccrual(
		cfg.AccuralSystemAddress,
		storages.OrdersRepo,
		storages.BalanceRepo,
		storages.Transactions,
	)

	//init server
	serverInstance := server.New(
		cfg.Address,
		storages.Transactions,
		storages.UsersRepo,
		storages.OrdersRepo,
		storages.BalanceRepo,
		storages.WithdrawalRepo,
		cfg.Secret,
	)

	//run accrual sync
	go accrualInstance.Sync()

	//run server
	go serverInstance.Start()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		serverInstance.Shutdown(context.Background(), &wg)
	}()

	go func() {
		storages.ShutDown(&wg)
	}()

	wg.Wait()
}
