package main

import (
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/config"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/orders"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/users"
	"log"
)

func main() {
	//init config
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	//init storage
	usersRepo, err := users.New(cfg.DBURI)
	if err != nil {
		log.Fatal(err)
	}
	ordersRepo, err := orders.New(cfg.DBURI)
	if err != nil {
		log.Fatal(err)
	}

	defer usersRepo.Close()
	defer ordersRepo.Close()

	//init server
	serverInstance := server.New(cfg.Address, usersRepo, ordersRepo, cfg.Secret)

	//run server
	if err = serverInstance.Start(); err != nil {
		log.Fatal(err)
	}
}
