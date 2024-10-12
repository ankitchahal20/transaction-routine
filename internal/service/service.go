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

// creditCardLimitOfferClient should only be created once throughtout the application lifetime
func NewTransactionRoutineService(conn db.TransactionRoutineService) *TransactionRoutineService {
	if transactionRoutineClient == nil {
		once.Do(
			func() {
				transactionRoutineClient = &TransactionRoutineService{
					repo: conn,
				}
			})
	} else {
		utils.Logger.Info("transaction routine client is alredy created")
	}
	return transactionRoutineClient
}
