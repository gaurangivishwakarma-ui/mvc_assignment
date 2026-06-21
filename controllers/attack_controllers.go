package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/gaurangi/mvc_assignment/models"
	"github.com/gaurangi/mvc_assignment/services"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AttackOpponent(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pgAttackerID := r.Context().Value(middleware.PlayerIDKey).(pgtype.UUID)

		var req models.InstantBattleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		result, statusCode, err := services.AttackOpponent(r.Context(), pool, pgAttackerID, req)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
