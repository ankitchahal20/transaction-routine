package db

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ankit/project/transaction-routine/internal/db/entities"
	"github.com/ankit/project/transaction-routine/internal/models"
	error "github.com/ankit/project/transaction-routine/internal/transactionroutineerror"
	"github.com/ankit/project/transaction-routine/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (p gormDB) CreateTransactions(ctx *gin.Context, transactionInfo models.Transactions, txid string) *error.TransactionRoutineError {
	// Use GORM's Create Transaction method
	createdTime := time.Now().Format(time.RFC3339)
	creationTime, parsedErr := time.Parse(time.RFC3339, createdTime)
	if parsedErr != nil {
		utils.Logger.Error("error while parsing creation time for creating transaction record", zap.String("transaction_id", txid))
		return &error.TransactionRoutineError{
			Trace:   txid,
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("unable to add transaction info : %v", fmt.Errorf("unable to parse creation time for creating transaction record")),
		}
	}

	transactionData := entities.Transactions{
		ID:              transactionInfo.TransactionID,
		AccountID:       transactionInfo.AccountID,
		OperationTypeID: transactionInfo.OperationTypeID,
		Amount:          *transactionInfo.Amount,
		EventDate:       creationTime,
	}

	tx := p.db.Begin()
	if err := tx.Create(&transactionData).Error; err != nil {
		utils.Logger.Error(fmt.Sprintf("error while inserting transaction details, txid: %v, error: %v", txid, err))

		// if transaction is already present
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &error.TransactionRoutineError{
				Trace:   txid,
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("transaction already present : %v", err.Error()),
			}
		}
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return &error.TransactionRoutineError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprint("Invalid account ID. The referenced account does not exist."),
				Trace:   txid,
			}
		}
		return &error.TransactionRoutineError{
			Trace:   txid,
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("unable to add transaction info : %v", err.Error()),
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Logger.Error("error while committing the transaction for create transaction ", zap.String("error : ", err.Error()), zap.String("transaction_id", txid))
		return &error.TransactionRoutineError{
			Trace:   txid,
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	utils.Logger.Info(fmt.Sprintf("successfully added the transaction entry in db, txid: %v", txid))
	return nil
}

func (p gormDB) GetTransaction(ctx *gin.Context, transactionID string, txid string) (models.Transactions, *error.TransactionRoutineError) {
	var scannedTransaction entities.Transactions
	var fetchedTransaction models.Transactions
	// Use GORM's First method to retrieve the account
	if err := p.db.Where("id = ?", transactionID).First(&scannedTransaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// case where no rows were found i.e. account not found
			return fetchedTransaction, &error.TransactionRoutineError{
				Code:    http.StatusNotFound,
				Message: "transaction not found",
				Trace:   txid,
			}
		}

		utils.Logger.Error(fmt.Sprintf("error while fetching transaction details from db, txid : %v, error: %v", txid, err))
		return fetchedTransaction, &error.TransactionRoutineError{
			Code:    http.StatusInternalServerError,
			Message: "unable to get the transactions details",
			Trace:   txid,
		}
	}
	fetchedTransaction.TransactionID = scannedTransaction.ID
	fetchedTransaction.AccountID = scannedTransaction.AccountID
	fetchedTransaction.Amount = &scannedTransaction.Amount
	fetchedTransaction.EventDate = scannedTransaction.EventDate.UTC()
	fetchedTransaction.OperationTypeID = scannedTransaction.OperationTypeID
	// Successfully fetched account
	utils.Logger.Info(fmt.Sprintf("successfully fetched transaction details from db, txid : %v", txid))
	return fetchedTransaction, nil
}
