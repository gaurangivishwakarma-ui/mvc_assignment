package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
)

func GetVillage(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pgPlayerID := r.Context().Value(middleware.PlayerIDKey).(pgtype.UUID)

		profile, err := queries.GetPlayerProfile(r.Context(), pgPlayerID)
		if err != nil {
			http.Error(w, "Failed to fetch player profile", http.StatusInternalServerError)
			return
		}

		buildings, err := queries.GetPlayerBuildings(r.Context(), pgPlayerID)
		if err != nil {
			buildings = []db.BuildingsOwned{}
		}

		response := map[string]interface{}{
			"profile":   profile,
			"buildings": buildings,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
