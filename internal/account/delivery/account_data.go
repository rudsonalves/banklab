package delivery

type AccountData struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Number     string `json:"number"`
	Branch     string `json:"branch"`
	Balance    int64  `json:"balance"`
	Status     string `json:"status"`
}

type TransferData struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	FromBalance   int64  `json:"from_balance"`
	ToBalance     int64  `json:"to_balance"`
}
