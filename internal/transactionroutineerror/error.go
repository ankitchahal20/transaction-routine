package transactionroutineerror

type TransactionRoutineError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}
