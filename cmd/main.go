package main

import (
	"net/http"
	"log"
	"github.com/go-chi/chi/v5"
	"encoding/json"
	"context"
    "fmt"
    "google.golang.org/genai"
)

type Recipe struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Ingredients []Item `json:"ingredients"`
}

type Item struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	Quantity int `json:"quantity"`
}

type Receipt struct {
	ID string `json:"id"`
	Store string `json:"store"`
	Status string `json:"status"`
	Items []Item `json:"items"`
}

func main() {
	// HTTP router
	router := chi.NewRouter()

	// Endpoint to check app health
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("App health good"))
	})

	// REST API
	router.Post("/receipts", createReceipt)
	router.Get("/receipts", getReceipts)
	router.Get("/receipts/{id}", getReceipt)
	router.Put("/receipts/{id}", updateReceipt)
	router.Delete("/receipts/{id}", deleteReceipt)

	// Run HTTP server
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)
}

// In progress: receiving raw image and saving receipt to database
// Plan: Add users after single-user functionality is complete

func createReceipt(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form 
    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

	// Receives image of receipt
	file, handler, err := r.FormFile("receipt_image")
    if err != nil {
        http.Error(w, "Error retrieving file", http.StatusBadRequest)
        return
    }
    defer file.Close()

	// Read file into mem so we can send it to Gemini
	buf := make([]byte, handler.Size)
    _, err = file.Read(buf)
    if err != nil {
        http.Error(w, "Error reading file", http.StatusInternalServerError)
        return
    }

	// Gemini reads receipt -> converts to json

	// Planned: Parse json to save to db
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	json.NewEncoder(w).Encode(receipt)
}

func getReceipts(w http.ResponseWriter, r *http.Request) {
	receipts := []Receipt{
		{
			ID: "1",
			Store: "Ralphs",
			Status: "Done",
			Items: []Item{
				{Name: "Tomatoes", Price: 0.99, Quantity:1},
				{Name: "Potatoes", Price: 0.50, Quantity:2},
				{Name: "Marshmellows", Price: 10.0, Quantity:5},
			},
		},
	}
	json.NewEncoder(w).Encode(receipts)
}

func getReceipt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	receipt := Receipt{
		ID:     id,
		Store:  "Trader Joe's",
		Status: "Done",
		Items: []Item{
			{Name: "Bread", Price: 2.49, Quantity: 2},
		},
	}
	json.NewEncoder(w).Encode(receipt)

}

func updateReceipt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	receipt := Receipt{
		ID:     id,
		Store:  "Updated Store",
		Status: "Updated",
		Items: []Item{
			{Name: "Updated Item", Price: 1.23, Quantity: 3},
		},
	}
	json.NewEncoder(w).Encode(receipt)
}

func deleteReceipt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res := map[string]string{
		"id":     id,
		"status": "deleted",
	}
	json.NewEncoder(w).Encode(res)
}