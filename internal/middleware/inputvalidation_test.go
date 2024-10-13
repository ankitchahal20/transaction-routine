package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankit/project/transaction-routine/internal/constants"
	"github.com/ankit/project/transaction-routine/internal/models"
	"github.com/ankit/project/transaction-routine/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateAccountRequestInput(t *testing.T) {
	// init logging client
	utils.InitLogClient()

	// Case 2 : documentNumber is missing
	requestFields := models.Accounts{}

	jsonValue, _ := json.Marshal(requestFields)

	w := httptest.NewRecorder()
	_, e := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewBuffer(jsonValue))
	req.Header.Add(constants.ContentType, "application/json")
	e.Use(ValidateInputRequest())
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateGetAccountRequestInput(t *testing.T) {
	// init logging client
	utils.InitLogClient()

	// case 1 : sending invalid uuid as account_id

	// Create a test context with the Gin engine
	r := gin.Default()
	r.Use(ValidateInputRequest())
	r.GET("/v1/accounts/:account_id", func(c *gin.Context) {})

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/accounts/", nil)
	ctx.Request.Header.Add(constants.ContentType, "application/json")
	ctx.Params = []gin.Param{
		{Key: "account_id", Value: ""},
	}

	// Serve the request through the Gin engine
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateTransactionPositiveAmount(t *testing.T) {
	// Initialize logging client
	utils.InitLogClient()

	r := gin.Default()
	r.Use(ValidateInputRequest())
	r.POST("/v1/transactions", func(c *gin.Context) {})
	amount := -100.10
	transactionInfo := models.Transactions{
		OperationTypeID: 4,
		AccountID:       "some-valid-uuid",
		Amount:          &amount, // Should be positive
	}

	body, _ := json.Marshal(transactionInfo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/transactions", bytes.NewReader(body))
	ctx.Request.Header.Add(constants.ContentType, "application/json")

	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateTransactionNegativeAmount(t *testing.T) {
	// Initialize logging client
	utils.InitLogClient()

	r := gin.Default()
	r.Use(ValidateInputRequest())
	r.POST("/v1/transactions", func(c *gin.Context) {})
	amount := -100.10
	transactionInfo := models.Transactions{
		OperationTypeID: 1,
		AccountID:       uuid.New().String(),
		Amount:          &amount,
	}

	body, _ := json.Marshal(transactionInfo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/transactions", bytes.NewReader(body))
	ctx.Request.Header.Add(constants.ContentType, "application/json")

	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestValidateCreateTransactionInput(t *testing.T) {
	// Initialize logging client
	utils.InitLogClient()

	// Create a test context with the Gin engine
	r := gin.Default()
	r.Use(ValidateInputRequest())
	r.POST("/v1/transactions", func(c *gin.Context) {})

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/transactions", nil)
	ctx.Request.Header.Add(constants.ContentType, "application/json")

	// Case 1: Missing operation_type_id
	ctx.Request.Body = httptest.NewRequest(http.MethodPost, "/v1/transactions", nil).Body
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	amount := -100.10
	// Case 2: Valid transaction creation
	transactionInfo := models.Transactions{
		OperationTypeID: 1,
		AccountID:       uuid.New().String(),
		Amount:          &amount,
	}
	body, _ := json.Marshal(transactionInfo)


	w = httptest.NewRecorder()
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/transactions", bytes.NewReader(body))
	r.ServeHTTP(w, ctx.Request)
	assert.Equal(t, http.StatusOK, w.Code)
}
