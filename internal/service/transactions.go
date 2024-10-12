package service

import (
	"fmt"
	"net/http"

	"github.com/ankit/project/transaction-routine/internal/constants"
	"github.com/ankit/project/transaction-routine/internal/models"
	"github.com/ankit/project/transaction-routine/internal/transactionroutineerror"
	"github.com/ankit/project/transaction-routine/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

func CreateTransaction() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.GetString(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for create account operation, txid : %v", txid))
		var transactionInfo models.Transactions
		if err := ctx.ShouldBindBodyWith(&transactionInfo, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("user request for account creation is unmarshalled successfully, txid : %v", txid))

			transactionID, err := transactionRoutineClient.createTransactions(ctx, transactionInfo, txid)
			if err != nil {
				utils.RespondWithError(ctx, err.Code, err.Message)
				return
			}

			// return 201 since new resource is created
			ctx.JSON(http.StatusCreated, map[string]string{
				"transaction_id": transactionID,
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *TransactionRoutineService) createTransactions(ctx *gin.Context, transactionInfo models.Transactions, txid string) (string, *transactionroutineerror.TransactionRoutineError) {
	// generate the accountID and customerID from uuid package and set in the the accountInfo
	transactionID := uuid.New().String()
	transactionInfo.TransactionID = transactionID

	utils.Logger.Info(fmt.Sprintf("calling db layer for account creation, txid : %v", txid))
	err := service.repo.CreateTransactions(ctx, transactionInfo, txid)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("received error from db layer during account creation, txid : %v", txid))
		return "", err
	}
	return transactionID, nil
}
