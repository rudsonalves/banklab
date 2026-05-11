package delivery

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

type AccountBalanceData struct {
	AccountID string `json:"account_id"`
	Balance   int64  `json:"balance"`
}

type InternalTransferRecipientsData struct {
	Accounts []InternalTransferRecipientData `json:"accounts"`
}

type InternalTransferRecipientData struct {
	AccountID     string `json:"account_id"`
	HolderName    string `json:"holder_name"`
	Document      string `json:"document"`
	Branch        string `json:"branch"`
	AccountNumber string `json:"account_number"`
}
