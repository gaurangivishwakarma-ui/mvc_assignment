package controllers

import (
	"encoding/json"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/gaurangi/mvc_assignment/services"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetMatch(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pgPlayerID := r.Context().Value(middleware.PlayerIDKey).(pgtype.UUID)

		result, statusCode, err := services.GetMatch(r.Context(), queries, pgPlayerID)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(result)
	}
}

func GetBattleReplay(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pgPlayerID := r.Context().Value(middleware.PlayerIDKey).(pgtype.UUID)

		battleIDStr := r.URL.Query().Get("id")
		if battleIDStr == "" {
			http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
			return
		}

		result, statusCode, err := services.GetBattleReplay(r.Context(), queries, pgPlayerID, battleIDStr)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
