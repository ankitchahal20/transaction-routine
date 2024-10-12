package models

import (
	"time"
)

type Accounts struct {
	AccountID      string `json:"account_id,omitempty"`
	DocumentNumber string `json:"document_number"`
}

type OperationType struct {
	OperationTypeID uint   `json:"operation_type_id"`
	Description     string `json:"description"`
}

type Transactions struct {
	TransactionID   string    `json:"transaction_id,omitempty"`
	AccountID       string    `json:"account_id"`
	OperationTypeID uint      `json:"operation_type_id"`
	Amount          *float64  `json:"amount"`
	EventDate       time.Time `json:"event_date,omitempty"`
}
