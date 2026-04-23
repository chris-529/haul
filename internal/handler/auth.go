package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chris-529/haul/internal/models"
	"github.com/jackc/pgx/v5/pgxpool" // Import the pool type
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	// Read in request as a user model
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash
	hashed, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	// Save to DB
	_, err := h.DB.Exec(context.Background(),
		"INSERT INTO users (email, password_hash) VALUES ($1, $2)",
		u.Email, string(hashed))

	// Check for registration error
	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	fmt.Println("Registration successful!")
	w.WriteHeader(http.StatusCreated)
}
