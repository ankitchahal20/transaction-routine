package db

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ankit/project/transaction-routine/internal/db/entities"
	"github.com/ankit/project/transaction-routine/internal/models"
	error "github.com/ankit/project/transaction-routine/internal/transactionroutineerror"
	"github.com/ankit/project/transaction-routine/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (p gormDB) CreateAccount(ctx *gin.Context, accountInfo models.Accounts, txid string) *error.TransactionRoutineError {
	// Use GORM's Create method
	accountData := entities.Accounts{
		ID:             accountInfo.AccountID,
		DocumentNumber: accountInfo.DocumentNumber,
	}

	tx := p.db.Begin()
	if err := tx.Create(&accountData).Error; err != nil {
		utils.Logger.Error(fmt.Sprintf("error while inserting account details, txid: %v, error: %v", txid, err))

		// if account is already present
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &error.TransactionRoutineError{
				Trace:   txid,
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("account already present : %v", err.Error()),
			}
		}
		return &error.TransactionRoutineError{
			Trace:   txid,
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("unable to add account info : %v", err.Error()),
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Logger.Error("error while committing the transaction for create account ", zap.String("error : ", err.Error()), zap.String("transaction_id", txid))
		return &error.TransactionRoutineError{
			Trace:   txid,
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	utils.Logger.Info(fmt.Sprintf("successfully added the account entry in db, txid: %v", txid))
	return nil
}

func (p gormDB) GetAccount(ctx *gin.Context, accountID string, txid string) (models.Accounts, *error.TransactionRoutineError) {
	var scannedAccount entities.Accounts
	var fetchedAccount models.Accounts
	// Use GORM's First method to retrieve the account
	if err := p.db.Where("id = ?", accountID).First(&scannedAccount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// case where no rows were found i.e. account not found
			return fetchedAccount, &error.TransactionRoutineError{
				Code:    http.StatusNotFound,
				Message: "account not found",
				Trace:   txid,
			}
		}

		utils.Logger.Error(fmt.Sprintf("error while fetching account from db, txid : %v, error: %v", txid, err))
		return fetchedAccount, &error.TransactionRoutineError{
			Code:    http.StatusInternalServerError,
			Message: "unable to get the account",
			Trace:   txid,
		}
	}
	fetchedAccount.AccountID = scannedAccount.ID
	fetchedAccount.DocumentNumber = scannedAccount.DocumentNumber
	// Successfully fetched account
	utils.Logger.Info(fmt.Sprintf("successfully fetched account from db, txid : %v", txid))
	return fetchedAccount, nil
}
