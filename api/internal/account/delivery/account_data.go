package delivery

import "time"

type AccountData struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Number     string `json:"number"`
	Branch     string `json:"branch"`
	Balance    int64  `json:"balance"`
	Status     string `json:"status"`
}

type AccountSummaryData struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Number     string `json:"number"`
	Branch     string `json:"branch"`
	Status     string `json:"status"`
}

type StatementItemData struct {
	TransactionID string    `json:"transaction_id"`
	Type          string    `json:"type"`
	Amount        int64     `json:"amount"`
	BalanceAfter  int64     `json:"balance_after"`
	ReferenceID   *string   `json:"reference_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type StatementData struct {
	AccountID  string               `json:"account_id"`
	Items      []StatementItemData  `json:"items"`
	NextCursor *StatementCursorData `json:"next_cursor"`
}

type StatementCursorData struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

type AccountBalanceData struct {
	AccountID string `json:"account_id"`
	Balance   int64  `json:"balance"`
}
