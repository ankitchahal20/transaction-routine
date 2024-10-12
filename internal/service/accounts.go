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

func CreateAccount() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.GetString(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for create account operation, txid : %v", txid))
		var accountInfo models.Accounts
		if err := ctx.ShouldBindBodyWith(&accountInfo, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("user request for account creation is unmarshalled successfully, txid : %v", txid))

			createdAccount, err := transactionRoutineClient.createAccount(ctx, accountInfo, txid)
			if err != nil {
				utils.RespondWithError(ctx, err.Code, err.Message)
				return
			}

			// return 201 since new resource is created
			ctx.JSON(http.StatusCreated, map[string]string{
				"account_id": createdAccount,
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *TransactionRoutineService) createAccount(ctx *gin.Context, accountInfo models.Accounts, txid string) (string, *transactionroutineerror.TransactionRoutineError) {
	// generate the accountID and customerID from uuid package and set in the the accountInfo
	accountID := uuid.New().String()
	accountInfo.AccountID = accountID

	utils.Logger.Info(fmt.Sprintf("calling db layer for account creation, txid : %v", txid))
	err := service.repo.CreateAccount(ctx, accountInfo, txid)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("received error from db layer during account creation, txid : %v", txid))
		return "", err
	}
	return accountID, nil
}

func GetAccount() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.GetString(constants.TransactionID)
		accountID := ctx.Param(constants.AccountID)
		utils.Logger.Info(fmt.Sprintf("request received for get %v account, txid : %v", accountID, txid))
		utils.Logger.Info(fmt.Sprintf("calling service layer for getting %v accountID, txid : %v", accountID, txid))
		fetchedAccount, err := transactionRoutineClient.getAccount(ctx, accountID, txid)
		if err != nil {
			utils.RespondWithError(ctx, err.Code, err.Message)
			return
		}

		ctx.JSON(http.StatusOK, fetchedAccount)
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

func (service *TransactionRoutineService) getAccount(ctx *gin.Context, accountID string, txid string) (models.Accounts, *transactionroutineerror.TransactionRoutineError) {
	utils.Logger.Info(fmt.Sprintf("calling db layer for getting %v account, txid : %v", accountID, txid))
	fetchedAccount, err := service.repo.GetAccount(ctx, accountID, txid)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("received error from db layer during getting %v account, txid : %v", accountID, txid))
		return models.Accounts{}, err
	}

	return fetchedAccount, nil
}
