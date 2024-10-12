package service

import (
	"sync"

	"github.com/ankit/project/transaction-routine/internal/db"
	"github.com/ankit/project/transaction-routine/internal/utils"
)

var (
	transactionRoutineClient *TransactionRoutineService
	once                     sync.Once
)

type TransactionRoutineService struct {
	repo db.TransactionRoutineService
}

func NewTransactionRoutineService(repo db.TransactionRoutineService) *TransactionRoutineService {
	if transactionRoutineClient == nil {
		once.Do(
			func() {
				transactionRoutineClient = &TransactionRoutineService{
					repo: repo,
				}
			})
	} else {
		utils.Logger.Info("transaction routine client is alredy created")
	}
	return transactionRoutineClient
}
