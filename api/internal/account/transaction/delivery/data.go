package delivery

type TransferData struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	FromBalance   int64  `json:"from_balance"`
	ToBalance     int64  `json:"to_balance"`
}
