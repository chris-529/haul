package main

import (
	"net/http"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
	"encoding/json"
	"context"
    "google.golang.org/genai"
)

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

var apiKey string

func main() {
	// Load API key from .env
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    apiKey = os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        log.Fatal("GEMINI_API_KEY not set")
    }

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
		writeJSONError(w, http.StatusBadRequest, "Error parsing form: "+err.Error())
		return
	}

	// Extract receipt image from request
	file, handler, err := r.FormFile("receipt_image")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Error retrieving file: "+err.Error())
		return
	}
	defer file.Close()

	// Read file into memory so we can send it to Gemini
	buf := make([]byte, handler.Size)
	_, err = file.Read(buf)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Error reading file: "+err.Error())
		return
	}

	// Extract the MIME type
	mimeType := handler.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	// Initialize Gemini Client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to initialize Gemini client: "+err.Error())
		return
	}

	// Only return JSON
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

	// Prompt with strict JSON structure
	prompt := `Analyze this receipt. Extract the store name and a list of all items purchased.
	You MUST respond with valid JSON matching this exact structure:
	{
	"id": "12345",
	"store": "Store Name Here",
	"status": "Done",
	"items": [
		{"name": "Item 1", "price": 1.99, "quantity": 1},
		{"name": "Item 2", "price": 5.50, "quantity": 2}
	]
	}`

	// Send image + prompt to Gemini
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					InlineData: &genai.Blob{
						MIMEType: mimeType,
						Data:     buf,
					},
				},
				{
					Text: prompt,
				},
			},
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		contents,
		config,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "AI generation failed: "+err.Error())
		return
	}

	// Unmarshal the JSON string directly into Go structs
	var receipt Receipt
	if err := json.Unmarshal([]byte(result.Text()), &receipt); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to parse AI response into structs: "+err.Error())
		return
	}

	// Return the JSON to the client
	// Plan: save receipt + items to database
	// Further planned: save receipts + items per user
	w.Header().Set("Content-Type", "application/json")
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

// Helper functions

func writeJSONError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{
        "error": msg,
    })
}