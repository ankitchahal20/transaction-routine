package transactionroutineerror

type TransactionRoutineError struct {
	Code    int    `json:"status_code"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}
