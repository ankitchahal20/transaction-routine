// Code generated by MockGen. DO NOT EDIT.
// Source: internal/db/db.go

// Package db is a generated GoMock package.
package db

import (
	reflect "reflect"

	models "github.com/ankit/project/transaction-routine/internal/models"
	transactionroutineerror "github.com/ankit/project/transaction-routine/internal/transactionroutineerror"
	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockTransactionRoutineService is a mock of TransactionRoutineService interface.
type MockTransactionRoutineService struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRoutineServiceMockRecorder
}

// MockTransactionRoutineServiceMockRecorder is the mock recorder for MockTransactionRoutineService.
type MockTransactionRoutineServiceMockRecorder struct {
	mock *MockTransactionRoutineService
}

// NewMockTransactionRoutineService creates a new mock instance.
func NewMockTransactionRoutineService(ctrl *gomock.Controller) *MockTransactionRoutineService {
	mock := &MockTransactionRoutineService{ctrl: ctrl}
	mock.recorder = &MockTransactionRoutineServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionRoutineService) EXPECT() *MockTransactionRoutineServiceMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockTransactionRoutineService) CreateAccount(arg0 *gin.Context, arg1 models.Accounts, arg2 string) *transactionroutineerror.TransactionRoutineError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0, arg1, arg2)
	ret0, _ := ret[0].(*transactionroutineerror.TransactionRoutineError)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockTransactionRoutineServiceMockRecorder) CreateAccount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockTransactionRoutineService)(nil).CreateAccount), arg0, arg1, arg2)
}

// CreateTransactions mocks base method.
func (m *MockTransactionRoutineService) CreateTransactions(arg0 *gin.Context, arg1 models.Transactions, arg2 string) *transactionroutineerror.TransactionRoutineError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransactions", arg0, arg1, arg2)
	ret0, _ := ret[0].(*transactionroutineerror.TransactionRoutineError)
	return ret0
}

// CreateTransactions indicates an expected call of CreateTransactions.
func (mr *MockTransactionRoutineServiceMockRecorder) CreateTransactions(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransactions", reflect.TypeOf((*MockTransactionRoutineService)(nil).CreateTransactions), arg0, arg1, arg2)
}

// GetAccount mocks base method.
func (m *MockTransactionRoutineService) GetAccount(arg0 *gin.Context, arg1, arg2 string) (models.Accounts, *transactionroutineerror.TransactionRoutineError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.Accounts)
	ret1, _ := ret[1].(*transactionroutineerror.TransactionRoutineError)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockTransactionRoutineServiceMockRecorder) GetAccount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockTransactionRoutineService)(nil).GetAccount), arg0, arg1, arg2)
}

// GetTransaction mocks base method.
func (m *MockTransactionRoutineService) GetTransaction(arg0 *gin.Context, arg1, arg2 string) (models.Transactions, *transactionroutineerror.TransactionRoutineError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.Transactions)
	ret1, _ := ret[1].(*transactionroutineerror.TransactionRoutineError)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockTransactionRoutineServiceMockRecorder) GetTransaction(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockTransactionRoutineService)(nil).GetTransaction), arg0, arg1, arg2)
}
