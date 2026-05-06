package delivery

type TransferData struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	FromBalance   int64  `json:"from_balance"`
	ToBalance     int64  `json:"to_balance"`
}

type NewTransferData struct {
	FromAccountBranch string `json:"from_branch"`
	FromAccountNumber string `json:"from_account_number"`
	ToAccountBranch   string `json:"to_branch"`
	ToAccountNumber   string `json:"to_account_number"`
	Amount            int64  `json:"amount"`
	FromBalance       int64  `json:"from_balance"`
	ToBalance         int64  `json:"to_balance"`
}
