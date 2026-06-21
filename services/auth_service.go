package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func RegisterPlayer(ctx context.Context, queries *db.Queries, req models.RegisterRequest) (map[string]interface{}, int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to process password")
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

	createdAuth, err := queries.CreatePlayerAuth(ctx, authParams)
	if err != nil {
		fmt.Printf("DB Error creating auth: %v\n", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to save player credentials")
	}

	_, err = queries.CreatePlayerProfile(ctx, pgPlayerID)
	if err != nil {
		fmt.Printf(" DB Error creating profile: %v\n", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to create player profile")
	}

	return map[string]interface{}{
		"message":   "Player registered successfully!",
		"player_id": createdAuth.PlayerID,
		"username":  createdAuth.Username,
	}, http.StatusCreated, nil
}

func LoginPlayer(ctx context.Context, queries *db.Queries, req models.LoginRequest) (map[string]interface{}, int, error) {
	authRecord, err := queries.GetAuthByUsername(ctx, req.Username)
	if err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("Invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(authRecord.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("Invalid username or password")
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
		return nil, http.StatusInternalServerError, fmt.Errorf("Could not generate token")
	}

	return map[string]interface{}{
		"message": "Welcome back to your village!",
		"token":   tokenString,
	}, http.StatusOK, nil
}
