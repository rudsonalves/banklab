package delivery

import "time"

type StatementItemData struct {
	TransactionID string    `json:"transaction_id"`
	Type          string    `json:"type"`
	Amount        int64     `json:"amount"`
	BalanceAfter  int64     `json:"balance_after"`
	ReferenceID   *string   `json:"reference_id"`
	Description   *string   `json:"description,omitempty"`
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
