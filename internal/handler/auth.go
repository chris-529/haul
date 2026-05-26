package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/chris-529/haul/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if u.Email == "" || u.Password == "" || u.InviteCode == "" {
		http.Error(w, "Email, password, and invite code are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	tx, err := h.DB.Begin(ctx)
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var inviteID string

	err = tx.QueryRow(ctx,
		`SELECT id
		 FROM invite_codes
		 WHERE code = $1
		   AND used_at IS NULL
		   AND expires_at > NOW()`,
		u.InviteCode,
	).Scan(&inviteID)

	if err != nil {
		http.Error(w, "Invalid or expired invite code", http.StatusUnauthorized)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES ($1, $2)`,
		u.Email,
		string(hashed),
	)
	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	_, err = tx.Exec(ctx,
		`UPDATE invite_codes
		 SET used_at = NOW()
		 WHERE id = $1`,
		inviteID,
	)
	if err != nil {
		http.Error(w, "Failed to use invite code", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	// Decode input
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Compare stored hash to input
	var userID string
	var storedHash string

	err := h.DB.QueryRow(context.Background(),
		"SELECT id, password_hash FROM users WHERE email = $1",
		input.Email,
	).Scan(&userID, &storedHash)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(input.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Create a jwt to return
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		http.Error(w, "JWT secret not configured", http.StatusInternalServerError)
		return
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   input.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
