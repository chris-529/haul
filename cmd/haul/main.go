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

	// Handler for Auth calls
	authH := &handler.AuthHandler{DB: db.Pool}

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// REST API
	r.Route("/receipts", func(r chi.Router) {
		r.Use(handler.AuthMiddleware)

		r.Post("/", h.CreateReceipt)
		r.Get("/", h.GetReceipts)
		r.Get("/{id}", h.GetReceipt)
		r.Put("/{id}", h.UpdateReceipt)
		r.Delete("/{id}", h.DeleteReceipt)
	})

	// Auth routes
	r.Post("/register", authH.Register)
	r.Post("/login", authH.Login)

	log.Println("Running on :8080")
	http.ListenAndServe(":8080", r)
}
