package main

import (
	"log"

	"github.com/ankit/project/transaction-routine/internal/config"
	"github.com/ankit/project/transaction-routine/internal/db"
	"github.com/ankit/project/transaction-routine/internal/server"
	"github.com/ankit/project/transaction-routine/internal/service"
	"github.com/ankit/project/transaction-routine/internal/utils"
	//_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Initializing the Log client
	utils.InitLogClient()

	// Initializing the GlobalConfig
	err := config.InitGlobalConfig()
	if err != nil {
		log.Fatalf("Unable to initialize global config")
	}

	// Establishing the connection to DB.
	postgres := db.Init()

	// Initializing the client for transaction routine service
	service.NewTransactionRoutineService(postgres)

	// Starting the server
	server.Start()
}
