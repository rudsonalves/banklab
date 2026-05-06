package delivery

type DepositRequest struct {
	Amount int64 `json:"amount"`
}

type WithdrawRequest struct {
	Amount int64 `json:"amount"`
}

type TransferRequest struct {
	FromAccountBranch string `json:"from_branch"`
	FromAccountNumber string `json:"from_account_number"`
	ToAccountBranch   string `json:"to_branch"`
	ToAccountNumber   string `json:"to_account_number"`
	Amount            int64  `json:"amount"`
	IdempotencyKey    string `json:"idempotency_key"`
}
