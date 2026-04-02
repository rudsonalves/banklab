package delivery

type AccountData struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Number     string `json:"number"`
	Branch     string `json:"branch"`
	Balance    int64  `json:"balance"`
	Status     string `json:"status"`
}
