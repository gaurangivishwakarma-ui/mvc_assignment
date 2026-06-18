package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func CollectResources(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		playerIDStr := r.Context().Value(middleware.PlayerIDKey).(string)
		parsedUUID, _ := uuid.Parse(playerIDStr)
		pgPlayerID := pgtype.UUID{Bytes: parsedUUID, Valid: true}

		buildings, err := queries.GetResourceBuildingsForCollection(r.Context(), pgPlayerID)
		if err != nil {
			http.Error(w, "Failed to fetch resource buildings", http.StatusInternalServerError)
			return
		}

		var totalGoldGenerated int32 = 0
		var totalElixirGenerated int32 = 0

		for _, b := range buildings {
			if !b.LastCollectedAt.Valid {
				continue
			}

			hoursPassed := time.Since(b.LastCollectedAt.Time).Hours()

			generated := int32(hoursPassed * b.ProductionRateFloat)

			if b.ResourceType == db.ResourceTypeGold {
				totalGoldGenerated += generated
			} else if b.ResourceType == db.ResourceTypeElixir {
				totalElixirGenerated += generated
			}
		}

		if totalGoldGenerated == 0 && totalElixirGenerated == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "No new resources to collect right now. Check back later!",
			})
			return
		}

		addParams := db.AddPlayerResourcesParams{
			GoldCoins: totalGoldGenerated,
			Elixir:    totalElixirGenerated,
			ID:        pgPlayerID,
		}

		updatedProfile, err := queries.AddPlayerResources(r.Context(), addParams)
		if err != nil {
			fmt.Printf("Error adding resources: %v\n", err)
			http.Error(w, "Failed to add resources to profile", http.StatusInternalServerError)
			return
		}

		err = queries.ResetCollectionTime(r.Context(), pgPlayerID)
		if err != nil {
			fmt.Printf("Error resetting timers: %v\n", err)
			http.Error(w, "Failed to reset collection timers", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Resources collected successfully!",
			"gold_looted":   totalGoldGenerated,
			"elixir_looted": totalElixirGenerated,
			"new_balance": map[string]int32{
				"gold":   updatedProfile.GoldCoins,
				"elixir": updatedProfile.Elixir,
			},
		})
	}
}
