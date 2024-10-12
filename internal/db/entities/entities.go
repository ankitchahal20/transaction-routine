package entities

import "time"

type Accounts struct {
	ID             string `gorm:"type:varchar(255);primaryKey"`
	DocumentNumber string `gorm:"type:varchar(255);not null"`
}

func (t *Accounts) TableName() string {
	/*
		// schema_name.service_name_table_name
		schema_name = pismodata
		service_name = transaction_routine
		table_name = accounts
	*/
	return "pismodata.transaction_routine_accounts_data"
}

type Transactions struct {
	ID              string    `gorm:"type:varchar(255);primaryKey"`            // Unique identifier for each transaction
	AccountID       string    `gorm:"type:varchar(255);not null"`              // Foreign key to the Account table (varchar for flexibility)
	OperationTypeID uint      `gorm:"type:int;not null"`                       // Foreign key to the OperationTypes table (int type)
	Amount          float64   `gorm:"type:numeric(10,2);not null"`             // Transaction amount with up to 10 digits and 2 decimal places
	EventDate       time.Time `gorm:"type:timestamptz;default:now();not null"` // Timestamp with timezone; default to current time
	Accounts        Accounts  `gorm:"foreignKey:AccountID;references:ID"`
}

func (t *Transactions) TableName() string {
	return "pismodata.transaction_routine_transactions_data"
}
