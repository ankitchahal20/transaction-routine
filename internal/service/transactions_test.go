package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankit/project/transaction-routine/internal/constants"
	"github.com/ankit/project/transaction-routine/internal/db"
	"github.com/ankit/project/transaction-routine/internal/models"
	"github.com/ankit/project/transaction-routine/internal/transactionroutineerror"
	"github.com/ankit/project/transaction-routine/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction_Success(t *testing.T) {
	ResetTransactionRoutineClient()
	utils.InitLogClient()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := db.NewMockTransactionRoutineService(ctrl)

	mockRepo.EXPECT().CreateTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	w := httptest.NewRecorder()
	amount := 100.0
	// Create valid transaction info
	transactionInfo := models.Transactions{
		Amount:          &amount,
		OperationTypeID: 1,
	}

	jsonBody, _ := json.Marshal(transactionInfo)
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constants.TransactionID, uuid.New().String())
	ctx.Request = req
	transactionRoutineClient = NewTransactionRoutineService(mockRepo)

	handler := CreateTransaction()
	handler(ctx)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
}

func TestCreateTransaction_Failure_BindingError(t *testing.T) {
	ResetTransactionRoutineClient()
	utils.InitLogClient()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/transactions", http.NoBody) // No body, which causes a binding error

	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constants.TransactionID, uuid.New().String())
	ctx.Request = req

	handler := CreateTransaction()
	handler(ctx)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTransaction_Success(t *testing.T) {
	ResetTransactionRoutineClient()
	utils.InitLogClient()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := db.NewMockTransactionRoutineService(ctrl)

	transactionID := uuid.New().String()
	amount := 100.0
	expectedTransaction := models.Transactions{
		TransactionID:   transactionID,
		Amount:          &amount,
		OperationTypeID: 2,
	}

	mockRepo.EXPECT().GetTransaction(gomock.Any(), transactionID, gomock.Any()).Return(expectedTransaction, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/transactions/"+transactionID, nil)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constants.TransactionID, transactionID)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: constants.TransactionId, Value: transactionID},
	}
	transactionRoutineClient = NewTransactionRoutineService(mockRepo)

	handler := GetTransaction()
	handler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the response matches the expected transaction data
	var fetchedTransaction models.Transactions
	err := json.Unmarshal(w.Body.Bytes(), &fetchedTransaction)
	assert.Nil(t, err)
	assert.Equal(t, expectedTransaction.TransactionID, fetchedTransaction.TransactionID)
}

func TestGetTransaction_Failure(t *testing.T) {
	ResetTransactionRoutineClient()
	utils.InitLogClient()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := db.NewMockTransactionRoutineService(ctrl)

	transactionID := uuid.New().String()
	mockError := transactionroutineerror.TransactionRoutineError{
		Code:    http.StatusInternalServerError,
		Message: "database error",
	}

	mockRepo.EXPECT().GetTransaction(gomock.Any(), transactionID, gomock.Any()).Return(models.Transactions{}, &mockError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/transactions/"+transactionID, nil)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constants.TransactionID, transactionID)
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: constants.TransactionId, Value: transactionID},
	}
	transactionRoutineClient = NewTransactionRoutineService(mockRepo)

	handler := GetTransaction()
	handler(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
