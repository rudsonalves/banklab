package delivery

type TransferData struct {
	FromAccountBranch    string `json:"from_branch"`
	FromAccountNumber    string `json:"from_account_number"`
	TransactionReference string `json:"transaction_reference"`
	ToAccountBranch      string `json:"to_branch"`
	ToAccountNumber      string `json:"to_account_number"`
	Amount               int64  `json:"amount"`
	FromBalance          int64  `json:"from_balance"`
	ToBalance            int64  `json:"to_balance"`
}

type TransferReceiptData struct {
	OperationType            string  `json:"operation_type"`
	Amount                   int64   `json:"amount"`
	Status                   string  `json:"status"`
	TransactionReference     string  `json:"transaction_reference"`
	OperationDate            string  `json:"operation_date"`
	SourceBranch             string  `json:"source_branch"`
	SourceAccountNumber      string  `json:"source_account_number"`
	DestinationBranch        string  `json:"destination_branch"`
	DestinationAccountNumber string  `json:"destination_account_number"`
	RecipientName            string  `json:"recipient_name"`
	Description              *string `json:"description,omitempty"`
}
