package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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

		pgPlayerID := r.Context().Value(middleware.PlayerIDKey).(pgtype.UUID)

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

func UpgradePlayerVillage(pool *pgxpool.Pool) http.HandlerFunc {
	queries := db.New(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		pgPlayerID := r.Context().Value(middleware.PlayerIDKey).(pgtype.UUID)

		tx, err := pool.Begin(r.Context())
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context())
		qtx := queries.WithTx(tx)

		player, err := qtx.GetPlayerProfile(r.Context(), pgPlayerID)
		if err != nil {
			http.Error(w, "Player profile not found", http.StatusNotFound)
			return
		}

		goldCost := player.VillageLevel * 10000

		if player.VillageLevel >= 4 {
			http.Error(w, "Your Village has already reached maximum level!", http.StatusBadRequest)
			return
		}

		if player.GoldCoins < goldCost {
			http.Error(w, "Insufficient Gold Coins to upgrade your Village!", http.StatusPaymentRequired)
			return
		}

		err = qtx.UpgradeVillageLevel(r.Context(), db.UpgradeVillageLevelParams{
			GoldCost: goldCost,
			PlayerID: pgPlayerID,
		})
		if err != nil {
			http.Error(w, "Failed to upgrade village level", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(r.Context()); err != nil {
			http.Error(w, "Failed to finalize village upgrade", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":           "Congratulations! Your Village level has increased!",
			"old_village_level": player.VillageLevel,
			"new_village_level": player.VillageLevel + 1,
			"gold_spent":        goldCost,
			"remaining_gold":    player.GoldCoins - goldCost,
		})
	}
}
