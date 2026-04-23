package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/chris-529/haul/internal/models"
	"github.com/go-chi/chi/v5"
	"google.golang.org/genai"
)

type ReceiptHandler struct {
	APIKey string
}

func (h *ReceiptHandler) CreateReceipt(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	file, handler, err := r.FormFile("receipt_image")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	buf := make([]byte, handler.Size)
	_, err = file.Read(buf)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	mimeType := handler.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  h.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		writeJSONError(w, 500, err.Error())
		return
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

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
				{Text: prompt},
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
		writeJSONError(w, 500, err.Error())
		return
	}

	var receipt models.Receipt
	json.Unmarshal([]byte(result.Text()), &receipt)

	json.NewEncoder(w).Encode(receipt)
}

func (h *ReceiptHandler) GetReceipts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]models.Receipt{})
}

func (h *ReceiptHandler) GetReceipt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *ReceiptHandler) UpdateReceipt(w http.ResponseWriter, r *http.Request) {}
func (h *ReceiptHandler) DeleteReceipt(w http.ResponseWriter, r *http.Request) {}

//Helper funcs

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
