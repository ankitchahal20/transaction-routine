package db

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ankit/project/transaction-routine/internal/models"
	"github.com/ankit/project/transaction-routine/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateTransaction_Success(t *testing.T) {
    // Create the mock GORM DB
    utils.InitLogClient()
    
    gormDB, mock := setupMockGorm(t)
	amount := 100.34
    transactionInfo := models.Transactions{
        TransactionID:   "123e4567-e89b-12d3-a456-426614174001",
        AccountID:       "123e4567-e89b-12d3-a456-426614174000",
        OperationTypeID: 1,
        Amount:          &amount,
    }

    mock.ExpectBegin()
    mock.ExpectQuery(`INSERT INTO "pismodata"."transaction_routine_transactions_data" \("id","account_id","operation_type_id","amount","event_date"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "event_date"`).
        WithArgs(transactionInfo.TransactionID, transactionInfo.AccountID, transactionInfo.OperationTypeID, transactionInfo.Amount, sqlmock.AnyArg()). // Use sqlmock.AnyArg() for EventDate
        WillReturnRows(sqlmock.NewRows([]string{"event_date"}).AddRow(time.Now()))
    mock.ExpectCommit()

    gdb := NewGormDB(gormDB)

    ctx := &gin.Context{}
    err := gdb.CreateTransactions(ctx, transactionInfo, "test-transaction-id")
    assert.Nil(t, err)

    mockErr := mock.ExpectationsWereMet()
    assert.NoError(t, mockErr)
}


func TestCreateTransaction_TransactionAlreadyExists(t *testing.T) {
    // Create the mock GORM DB
    utils.InitLogClient()

    gormDB, mock := setupMockGorm(t)
	amount := 100.34
    transactionInfo := models.Transactions{
        TransactionID: "123e4567-e89b-12d3-a456-426614174001",
        Amount:        &amount,
        AccountID:     "123e4567-e89b-12d3-a456-426614174000",
    }

    mock.ExpectBegin()
    mock.ExpectQuery(`INSERT INTO "pismodata"."transaction_routine_transactions_data"`).
        WithArgs(transactionInfo.TransactionID, transactionInfo.AccountID, transactionInfo.OperationTypeID,transactionInfo.Amount, sqlmock.AnyArg()).
        WillReturnError(gorm.ErrDuplicatedKey)
    mock.ExpectRollback()

    gdb := NewGormDB(gormDB)

    ctx := &gin.Context{}
    err := gdb.CreateTransactions(ctx, transactionInfo, "test-transaction-id")
	fmt.Println("err : ", err)
    assert.NotNil(t, err)
    assert.Equal(t, http.StatusBadRequest, err.Code)

    mockErr := mock.ExpectationsWereMet()
    assert.Error(t, mockErr)
}

func TestGetTransaction_Success(t *testing.T) {
    gormDB, mock := setupMockGorm(t)
    utils.InitLogClient()

    transactionID := "123e4567-e89b-12d3-a456-426614174001"

    mock.ExpectQuery(`SELECT \* FROM "pismodata"."transaction_routine_transactions_data" WHERE id = \$1 ORDER BY "transaction_routine_transactions_data"."id" LIMIT \$2`).
        WithArgs(transactionID, 1).
        WillReturnRows(sqlmock.NewRows([]string{"id", "amount", "account_id"}).
            AddRow(transactionID, 1000, "123e4567-e89b-12d3-a456-426614174000"))

    gdb := NewGormDB(gormDB)

    ctx := &gin.Context{}

    fetchedTransaction, err := gdb.GetTransaction(ctx, transactionID, "test-transaction-id")
    assert.Nil(t, err)
    assert.NotNil(t, fetchedTransaction)

    mockErr := mock.ExpectationsWereMet()
    assert.NoError(t, mockErr)
}

func TestGetTransaction_NotFound(t *testing.T) {
    utils.InitLogClient()
    gormDB, mock := setupMockGorm(t)

    transactionID := "123e4567-e89b-12d3-a456-426614174001"

    mock.ExpectQuery(`SELECT \* FROM "pismodata"."transaction_routine_transactions_data" WHERE id = \$1 ORDER BY "transaction_routine_transactions_data"."id" LIMIT \$2`).
        WithArgs(transactionID, 1).
        WillReturnError(gorm.ErrRecordNotFound)

    gdb := NewGormDB(gormDB)
    ctx := &gin.Context{}

    _, err := gdb.GetTransaction(ctx, transactionID, "test-transaction-id")
    assert.NotNil(t, err)
    assert.Equal(t, http.StatusNotFound, err.Code)

    mockErr := mock.ExpectationsWereMet()
    assert.NoError(t, mockErr)
}

func TestGetTransaction_Error(t *testing.T) {
    utils.InitLogClient()

    gormDB, mock := setupMockGorm(t)

    transactionID := "123e4567-e89b-12d3-a456-426614174001"

    mock.ExpectQuery(`SELECT \* FROM "pismodata"."transaction_routine_transactions_data" WHERE id = \$1 ORDER BY "transaction_routine_transactions_data"."id" LIMIT \$2`).
        WithArgs(transactionID, 1).
        WillReturnError(fmt.Errorf("some error"))

    gdb := NewGormDB(gormDB)

    ctx := &gin.Context{}

    fetchedTransaction, err := gdb.GetTransaction(ctx, transactionID, "test-transaction-id")
    assert.NotNil(t, err)
    assert.Equal(t, http.StatusInternalServerError, err.Code)
    assert.Empty(t, fetchedTransaction.TransactionID)

    mockErr := mock.ExpectationsWereMet()
    assert.NoError(t, mockErr)
}
