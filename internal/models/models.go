package models

import "time"

type Item struct {
	ID        string  `json:"id"`
	ReceiptID string  `json:"receipt_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Unit      string  `json:"unit"`
}

type Receipt struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Store     string    `json:"store"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Items     []Item    `json:"items"`
}

type Recipe struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Ingredients []Item `json:"ingredients"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}
