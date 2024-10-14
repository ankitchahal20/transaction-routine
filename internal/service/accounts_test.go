package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
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

func ResetTransactionRoutineClient() {
	once = sync.Once{}
	transactionRoutineClient = nil
}

func TestCreateAccount_Success(t *testing.T) {
	ResetTransactionRoutineClient()
	utils.InitLogClient()
	// Setup the mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the repository
	mockRepo := db.NewMockTransactionRoutineService(ctrl)

	// Expected behavior for the mock repository
	mockRepo.EXPECT().CreateAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	w := httptest.NewRecorder()

	// Create a valid account info
	accountInfo := models.Accounts{
		DocumentNumber: "123e4567-e89b-12d3-a456-426614174000",
	}

	// Marshal the account info to JSON
	jsonBody, _ := json.Marshal(accountInfo)
	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constants.TransactionID, uuid.New().String())
	ctx.Request = req
	transactionRoutineClient = NewTransactionRoutineService(mockRepo)
	// Call the CreateAccount
	handler := CreateAccount()
	handler(ctx)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetAccount_Success(t *testing.T) {
	ResetTransactionRoutineClient()
	utils.InitLogClient()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := db.NewMockTransactionRoutineService(ctrl)

	// Mock data: Define the expected account data to be returned
	accountID := "123e4567-e89b-12d3-a456-426614174000"
	expectedAccount := models.Accounts{
		DocumentNumber: "123e4567-e89b-12d3-a456-426614174000",
		AccountID:      "123e4567-e89b-12d3-a456-426614174001",
	}

	mockRepo.EXPECT().GetAccount(gomock.Any(), accountID, gomock.Any()).Return(expectedAccount, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/account/"+accountID, nil)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(constants.TransactionID, uuid.New().String())
	ctx.Request = req
	ctx.Params = gin.Params{
		gin.Param{Key: "account_id", Value: accountID},
	}
	transactionRoutineClient = NewTransactionRoutineService(mockRepo)
	handler := GetAccount()
	handler(ctx)

	// Check for the proper status code and response
	assert.Equal(t, http.StatusOK, w.Code)

	// You can also check if the response matches the expected account data
	var fetchedAccount models.Accounts
	err := json.Unmarshal(w.Body.Bytes(), &fetchedAccount)
	assert.Nil(t, err)
	assert.Equal(t, expectedAccount.AccountID, fetchedAccount.AccountID)
}

func TestCreateAccount_Failure(t *testing.T) {
	ResetTransactionRoutineClient()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := db.NewMockTransactionRoutineService(ctrl)
	accountInfo := models.Accounts{
		DocumentNumber: uuid.New().String(),
	}

	// Mock an error response from the database layer
	mockError := transactionroutineerror.TransactionRoutineError{
		Code:    http.StatusInternalServerError,
		Message: "database error",
	}

	mockRepo.EXPECT().CreateAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mockError)

	jsonBody, _ := json.Marshal(accountInfo)
	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set(constants.TransactionID, uuid.New().String())
	ctx.Request = req
	transactionRoutineClient = NewTransactionRoutineService(mockRepo)

	// Call CreateAccount
	handler := CreateAccount()
	handler(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
}
