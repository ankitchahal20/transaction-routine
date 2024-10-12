package utils

import (
	"github.com/ankit/project/transaction-routine/internal/constants"
	"github.com/ankit/project/transaction-routine/internal/transactionroutineerror"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogClient() {
	Logger, _ = zap.NewDevelopment()
}

func RespondWithError(c *gin.Context, statusCode int, message string) {
	txid, _ := c.Get(constants.TransactionID)
	c.AbortWithStatusJSON(statusCode, transactionroutineerror.TransactionRoutineError{
		Trace:   txid.(string),
		Code:    statusCode,
		Message: message,
	})
}
