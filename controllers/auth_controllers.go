package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterPlayer(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON request", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to process password", http.StatusInternalServerError)
			return
		}

		playerID := uuid.New()

		pgPlayerID := pgtype.UUID{
			Bytes: playerID,
			Valid: true,
		}

		authParams := db.CreatePlayerAuthParams{
			PlayerID:     pgPlayerID,
			Username:     req.Username,
			PasswordHash: string(hashedPassword),
		}

		createdAuth, err := queries.CreatePlayerAuth(r.Context(), authParams)
		if err != nil {
			fmt.Printf("DB Error creating auth: %v\n", err)
			http.Error(w, "Failed to save player credentials", http.StatusInternalServerError)
			return
		}

		_, err = queries.CreatePlayerProfile(r.Context(), pgPlayerID)
		if err != nil {
			fmt.Printf(" DB Error creating profile: %v\n", err)
			http.Error(w, "Failed to create player profile", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Player registered successfully!",
			"player_id": createdAuth.PlayerID,
			"username":  createdAuth.Username,
		})
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginPlayer(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON request", http.StatusBadRequest)
			return
		}

		authRecord, err := queries.GetAuthByUsername(r.Context(), req.Username)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(authRecord.PasswordHash), []byte(req.Password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		playerIDStr := uuid.UUID(authRecord.PlayerID.Bytes).String()

		claims := jwt.MapClaims{
			"player_id": playerIDStr,
			"username":  authRecord.Username,
			"exp":       time.Now().Add(time.Hour * 72).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		secretKey := os.Getenv("JWT_SECRET")
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			fmt.Printf("Error signing token: %v\n", err)
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Welcome back to your village!",
			"token":   tokenString,
		})
	}
}
