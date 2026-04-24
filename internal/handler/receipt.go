package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/chris-529/haul/internal/db"
	"github.com/chris-529/haul/internal/models"
	"github.com/go-chi/chi/v5"
	"google.golang.org/genai"
)

type ReceiptHandler struct {
	APIKey string
}

func (h *ReceiptHandler) CreateReceipt(w http.ResponseWriter, r *http.Request) {

	// Get userID from JWT
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "Missing user ID")
		return
	}

	// Limit request body to 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	// Parse form for file
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	file, _, err := r.FormFile("receipt_image")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	// Read in image to memory
	buf, err := io.ReadAll(file)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Detect file type
	mimeType := detectImageMIME(buf)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedTypes[mimeType] {
		writeJSONError(w, http.StatusBadRequest, "unsupported file type")
		return
	}

	// Init Gemini cli
	ctx := r.Context()
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
	"store": "Store Name Here",
	"status": "Done",
	"items": [
		{"name": "Item 1", "price": 1.99, "quantity": 1, "unit": ""},
		{"name": "Item 2", "price": 5.50, "quantity": 2, "unit": ""}
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

	// Generate a response
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

	// Parse the gemini response into a Receipt model
	var receipt models.Receipt
	if err := json.Unmarshal([]byte(result.Text()), &receipt); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to parse AI response")
		return
	}

	// Make sure it has items
	if len(receipt.Items) == 0 {
		writeJSONError(w, http.StatusBadRequest, "no receipt items detected")
		return
	}

	receipt.UserID = userID
	receipt.Status = "Done"

	// Save this receipt to the database for userID
	if err := db.SaveReceipt(ctx, userID, &receipt); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// For now just return the receipt json back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}

func (h *ReceiptHandler) GetReceipts(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "Missing user ID")
		return
	}

	receipts, err := db.GetReceipts(r.Context(), userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipts)
}

func (h *ReceiptHandler) GetReceipt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *ReceiptHandler) UpdateReceipt(w http.ResponseWriter, r *http.Request) {}
func (h *ReceiptHandler) DeleteReceipt(w http.ResponseWriter, r *http.Request) {}

//Helper funcs

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}

func detectImageMIME(buf []byte) string {
	mimeType := http.DetectContentType(buf)

	if mimeType == "application/octet-stream" && len(buf) >= 12 {
		if string(buf[0:4]) == "RIFF" && string(buf[8:12]) == "WEBP" {
			return "image/webp"
		}
	}

	return mimeType
}
