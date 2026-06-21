package controllers

import (
	"encoding/json"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/gaurangi/mvc_assignment/models"
	"github.com/gaurangi/mvc_assignment/services"
	"github.com/jackc/pgx/v5/pgtype"
)

func PurchaseBuilding(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.PurchaseBuildingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON request", http.StatusBadRequest)
			return
		}

		pgPlayerID := r.Context().Value(middleware.PlayerIDKey).(pgtype.UUID)

		result, statusCode, err := services.PurchaseBuilding(r.Context(), queries, pgPlayerID, req)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(result)
	}
}

func CompleteBuild(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pgPlayerID := r.Context().Value(middleware.PlayerIDKey).(pgtype.UUID)

		var req models.CompleteBuildRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON request", http.StatusBadRequest)
			return
		}

		result, statusCode, err := services.CompleteBuild(r.Context(), queries, pgPlayerID, req)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(result)
	}
}
