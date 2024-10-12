package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ankit/project/transaction-routine/internal/constants"
	"github.com/ankit/project/transaction-routine/internal/models"
	"github.com/ankit/project/transaction-routine/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// This function gets the unique transactionID
func getTransactionID(c *gin.Context) string {
	transactionID := c.GetHeader(constants.TransactionID)
	_, err := uuid.Parse(transactionID)
	if err != nil {
		transactionID = uuid.New().String()
		c.Set(constants.TransactionID, transactionID)
	}
	return transactionID
}

func ValidateInputRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get the transactionID from headers if not present create a new.
		transactionID := getTransactionID(ctx)
		path := ctx.Request.URL.String()
		switch {
		case strings.Contains(path, constants.Accounts) && ctx.Request.Method == http.MethodPost:
			validateCreateAccountInput(ctx, transactionID)
		case strings.Contains(path, constants.Accounts) && ctx.Request.Method == http.MethodGet:
			validateGetAccountInput(ctx, transactionID)
		case strings.Contains(path, constants.Transactions) && ctx.Request.Method == http.MethodPost:
			validateCreateTransactionInput(ctx, transactionID)
		}

		ctx.Next()
	}
}

func validateCreateAccountInput(ctx *gin.Context, txid string) {
	var accountInfo models.Accounts
	err := ctx.ShouldBindBodyWith(&accountInfo, binding.JSON)
	if err != nil {
		utils.Logger.Error("error while unmarshaling the request field for create account data validation")
		utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidBodyCreateAccount)
		return
	}

	if accountInfo.DocumentNumber == "" {
		utils.Logger.Error(fmt.Sprintf("document number field is missing while creating an account, txid : %v", txid))
		errMessage := "document_number field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
}

func validateGetAccountInput(ctx *gin.Context, txid string) {
	accountID := ctx.Param(constants.AccountID)
	utils.Logger.Info(fmt.Sprintf("request received for get %v account, txid : %v", accountID, txid))
	_, erraccountUUID := uuid.Parse(accountID)
	if erraccountUUID != nil {
		utils.Logger.Error(fmt.Sprintf("Error parsing the %v accountID, txid : %v", accountID, txid))
		utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidAccountID)
		return
	}
}

func validateCreateTransactionInput(ctx *gin.Context, txid string) {
	var transactionInfo models.Transaction
	err := ctx.ShouldBindBodyWith(&transactionInfo, binding.JSON)
	if err != nil {
		utils.Logger.Error("error while unmarshaling the request field for create transaction")
		utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidBodyCreateTransaction)
		return
	}

	if transactionInfo.OperationTypeID == 0 {
		utils.Logger.Error(fmt.Sprintf("operation type id field is missing while creating a transaction, txid : %v", txid))
		errMessage := "operation_type_id field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
	if transactionInfo.AccountID == "" {
		utils.Logger.Error(fmt.Sprintf("account_id field is missing while creating a transaction, txid : %v", txid))
		errMessage := "account_id field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
	if transactionInfo.Amount == nil {
		utils.Logger.Error(fmt.Sprintf("amount field is missing while creating a transaction, txid : %v", txid))
		errMessage := "amount field is missing"
		utils.RespondWithError(ctx, http.StatusBadRequest, errMessage)
		return
	}
}
