package models

type Item struct {
	ID        string  `json:"id"`
	ReceiptID string  `json:"receipt_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Unit      string  `json:"unit"`
}

type Receipt struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Store  string `json:"store"`
	Status string `json:"status"`
	Items  []Item `json:"items"`
}

type Recipe struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Ingredients []Item `json:"ingredients"`
}
