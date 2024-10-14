package db

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ankit/project/transaction-routine/internal/models"
	"github.com/ankit/project/transaction-routine/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockGorm(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	utils.InitLogClient()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock db: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error initializing gorm db: %v", err)
	}

	return gormDB, mock
}

func TestCreateAccount_Success(t *testing.T) {
	// Create the mock GORM DB
	utils.InitLogClient()

	gormDB, mock := setupMockGorm(t)

	accountInfo := models.Accounts{
		AccountID:      "123e4567-e89b-12d3-a456-426614174000",
		DocumentNumber: "123456789",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "pismodata"."transaction_routine_accounts_data"`).
		WithArgs(accountInfo.AccountID, accountInfo.DocumentNumber).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gdb := NewGormDB(gormDB)

	ctx := &gin.Context{}
	err := gdb.CreateAccount(ctx, accountInfo, "test-transaction-id")
	assert.Nil(t, err)

	mockErr := mock.ExpectationsWereMet()
	assert.NoError(t, mockErr)
}

func TestCreateAccount_AccountAlreadyExists(t *testing.T) {
	// Create the mock GORM DB
	utils.InitLogClient()

	gormDB, mock := setupMockGorm(t)

	accountInfo := models.Accounts{
		AccountID:      "123e4567-e89b-12d3-a456-426614174000",
		DocumentNumber: "123456789",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "pismodata"."transaction_routine_accounts_data"`).
		WithArgs(accountInfo.AccountID, accountInfo.DocumentNumber).
		WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	gdb := NewGormDB(gormDB)

	ctx := &gin.Context{}

	err := gdb.CreateAccount(ctx, accountInfo, "test-transaction-id")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)

	mockErr := mock.ExpectationsWereMet()
	assert.Error(t, mockErr)
}

func TestGetAccount_Success(t *testing.T) {
	gormDB, mock := setupMockGorm(t)
	utils.InitLogClient()
	accountID := "123e4567-e89b-12d3-a456-426614174000"

	mock.ExpectQuery(`SELECT \* FROM "pismodata"."transaction_routine_accounts_data" WHERE id = \$1 ORDER BY "transaction_routine_accounts_data"."id" LIMIT \$2`).
		WithArgs(accountID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "document_number"}).
			AddRow(accountID, "123456789"))

	gdb := NewGormDB(gormDB)

	ctx := &gin.Context{}

	fetchedAccount, err := gdb.GetAccount(ctx, accountID, "test-transaction-id")
	assert.Nil(t, err)
	assert.NotNil(t, fetchedAccount)

	mockErr := mock.ExpectationsWereMet()
	assert.NoError(t, mockErr)
}

func TestGetAccount_NotFound(t *testing.T) {
	utils.InitLogClient()
	gormDB, mock := setupMockGorm(t)

	accountID := "123e4567-e89b-12d3-a456-426614174000"

	mock.ExpectQuery(`SELECT \* FROM "pismodata"."transaction_routine_accounts_data" WHERE id = \$1 ORDER BY "transaction_routine_accounts_data"."id" LIMIT \$2`).
		WithArgs(accountID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	gdb := NewGormDB(gormDB)
	ctx := &gin.Context{}

	_, err := gdb.GetAccount(ctx, accountID, "test-transaction-id")
	fmt.Println("err : ", err)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Code)

	mockErr := mock.ExpectationsWereMet()
	assert.NoError(t, mockErr)
}

func TestGetAccount_Error(t *testing.T) {
	utils.InitLogClient()

	gormDB, mock := setupMockGorm(t)

	accountID := "123e4567-e89b-12d3-a456-426614174000"

	mock.ExpectQuery(`SELECT \* FROM "pismodata"."transaction_routine_accounts_data" WHERE id = \$1 ORDER BY "transaction_routine_accounts_data"."id" LIMIT \$2`).
		WithArgs(accountID, 1).
		WillReturnError(fmt.Errorf("some error"))

	gdb := NewGormDB(gormDB)

	ctx := &gin.Context{}

	fetchedAccount, err := gdb.GetAccount(ctx, accountID, "test-transaction-id")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Empty(t, fetchedAccount.AccountID)

	mockErr := mock.ExpectationsWereMet()
	assert.NoError(t, mockErr)
}
