package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpgradeRequest struct {
	PlacementID string `json:"placement_id"`
}

func UpgradeBuilding(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		playerIDStr := r.Context().Value(middleware.PlayerIDKey).(string)
		parsedPlayerUUID, _ := uuid.Parse(playerIDStr)
		pgPlayerID := pgtype.UUID{Bytes: parsedPlayerUUID, Valid: true}

		var req UpgradeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		parsedPlacementUUID, err := uuid.Parse(req.PlacementID)
		if err != nil {
			http.Error(w, "Invalid placement ID format", http.StatusBadRequest)
			return
		}
		pgPlacementID := pgtype.UUID{Bytes: parsedPlacementUUID, Valid: true}

		ownedArgs := db.GetOwnedBuildingDetailsParams{
			ID:       pgPlacementID,
			PlayerID: pgPlayerID,
		}
		ownedBuilding, err := queries.GetOwnedBuildingDetails(r.Context(), ownedArgs)
		if err != nil {
			http.Error(w, "Building not found or you don't own it", http.StatusNotFound)
			return
		}

		nextLevelArgs := db.GetNextLevelBuildingParams{
			BuildingType: ownedBuilding.BuildingType,
			Level:        ownedBuilding.CurrentLevel + 1,
		}
		nextLevel, err := queries.GetNextLevelBuilding(r.Context(), nextLevelArgs)
		if err != nil {
			http.Error(w, "Building is already at max level!", http.StatusBadRequest)
			return
		}

		profile, err := queries.GetPlayerProfile(r.Context(), pgPlayerID)
		if err != nil || profile.VillageLevel < nextLevel.LevelReq {
			http.Error(w, "Your Village Level is too low to upgrade this building", http.StatusForbidden)
			return
		}

		var goldCost, elixirCost int32 = 0, 0
		if strings.ToLower(string(nextLevel.CostType)) == "gold" {
			goldCost = nextLevel.BuildCost
		} else if strings.ToLower(string(nextLevel.CostType)) == "elixir" {
			elixirCost = nextLevel.BuildCost
		}

		payArgs := db.PayForUpgradeParams{
			GoldCoins: goldCost,
			Elixir:    elixirCost,
			ID:        pgPlayerID,
		}
		if err := queries.PayForUpgrade(r.Context(), payArgs); err != nil {
			http.Error(w, "Insufficient resources to upgrade", http.StatusPaymentRequired)
			return
		}

		upgradeArgs := db.CommitBuildingUpgradeParams{
			BuildingID:   nextLevel.NextBuildingID,
			CurrentLevel: ownedBuilding.CurrentLevel + 1,
			ID:           pgPlacementID,
		}
		if err := queries.CommitBuildingUpgrade(r.Context(), upgradeArgs); err != nil {
			http.Error(w, "Failed to apply upgrade to building", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Building upgraded successfully!",
			"new_level":     ownedBuilding.CurrentLevel + 1,
			"building_name": nextLevel.Name,
		})
	}
}
