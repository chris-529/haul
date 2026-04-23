package main

import (
	"log"
	"net/http"
	"os"

	"github.com/chris-529/haul/internal/db"
	"github.com/chris-529/haul/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db.Connect()
	defer db.Close()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing API key")
	}

	// HTTP router
	r := chi.NewRouter()

	// Handler for API calls
	h := &handler.ReceiptHandler{
		APIKey: apiKey,
	}

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// REST API
	r.Route("/receipts", func(r chi.Router) {
		r.Post("/", h.CreateReceipt)
		r.Get("/", h.GetReceipts)
		r.Get("/{id}", h.GetReceipt)
		r.Put("/{id}", h.UpdateReceipt)
		r.Delete("/{id}", h.DeleteReceipt)
	})

	log.Println("Running on :8080")
	http.ListenAndServe(":8080", r)
}
