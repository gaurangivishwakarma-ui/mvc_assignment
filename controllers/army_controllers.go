package controllers

import (
	"encoding/json"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	_ "github.com/gaurangi/mvc_assignment/middleware"
)

func GetArmyCatalog(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		catalog, err := queries.GetTroopCatalog(r.Context())
		if err != nil {
			http.Error(w, "Failed to fetch troop catalog", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"troops_available": catalog,
		})
	}
}
