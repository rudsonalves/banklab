package delivery

type CreateAccountRequest struct {
	CustomerID string `json:"customer_id"`
}

type DepositRequest struct {
	Amount int64 `json:"amount"`
}

type WithdrawRequest struct {
	Amount int64 `json:"amount"`
}
