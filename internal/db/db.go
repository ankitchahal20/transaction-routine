package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/ankit/project/transaction-routine/internal/config"
	"github.com/ankit/project/transaction-routine/internal/constants"
	"github.com/ankit/project/transaction-routine/internal/db/entities"
	"github.com/ankit/project/transaction-routine/internal/models"
	"github.com/ankit/project/transaction-routine/internal/transactionroutineerror"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	conn *gorm.DB
	once sync.Once
)

type gormDB struct{ db *gorm.DB }

type TransactionRoutineService interface {
	CreateAccount(*gin.Context, models.Accounts, string) *transactionroutineerror.TransactionRoutineError
	GetAccount(*gin.Context, string, string) (models.Accounts, *transactionroutineerror.TransactionRoutineError)
	CreateTransactions(*gin.Context, models.Transactions, string) *transactionroutineerror.TransactionRoutineError
}

func Init() gormDB {
	once.Do(func() {
		cfg := config.GetConfig()
		connString := fmt.Sprintf(
			"host=%s dbname=%s password=%s user=%s port=%d, search_path=%s",
			cfg.Database.Host, cfg.Database.DBname, cfg.Database.Password,
			cfg.Database.User, cfg.Database.Port, constants.DBSchemaName,
		)
		var err error
		conn, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}

		// Create the schema if it doesn't exist
		err = conn.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", constants.DBSchemaName)).Error
		if err != nil {
			log.Fatalln("Failed to create schema:", err)
		}

		conn.AutoMigrate(&entities.Accounts{}, entities.Transactions{})
	})

	return gormDB{
		db: conn,
	}
}
