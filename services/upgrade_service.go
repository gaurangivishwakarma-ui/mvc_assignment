package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpgradeBuilding(ctx context.Context, queries *db.Queries, pgPlayerID pgtype.UUID, req models.UpgradeRequest) (map[string]interface{}, int, error) {
	parsedPlacementUUID, err := uuid.Parse(req.PlacementID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid placement ID format")
	}
	pgPlacementID := pgtype.UUID{Bytes: parsedPlacementUUID, Valid: true}
	_ = queries.AutoCompleteBuildings(ctx, pgPlayerID)

	ownedArgs := db.GetOwnedBuildingDetailsParams{
		ID:       pgPlacementID,
		PlayerID: pgPlayerID,
	}
	ownedBuilding, err := queries.GetOwnedBuildingDetails(ctx, ownedArgs)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Building not found or you don't own it")
	}

	if !ownedBuilding.IsBuilt.Bool {
		return nil, http.StatusBadRequest, fmt.Errorf("Cannot upgrade building: 5 seconds must pass for creating or upgrading it")
	}

	nextLevelArgs := db.GetNextLevelBuildingParams{
		BuildingType: ownedBuilding.BuildingType,
		Level:        ownedBuilding.CurrentLevel + 1,
	}
	nextLevel, err := queries.GetNextLevelBuilding(ctx, nextLevelArgs)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Building is already at max level!")
	}

	profile, err := queries.GetPlayerProfile(ctx, pgPlayerID)
	if err != nil || profile.VillageLevel < nextLevel.LevelReq {
		return nil, http.StatusForbidden, fmt.Errorf("Your Village Level is too low to upgrade this building")
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
	if err := queries.PayForUpgrade(ctx, payArgs); err != nil {
		return nil, http.StatusPaymentRequired, fmt.Errorf("Insufficient resources to upgrade")
	}

	upgradeArgs := db.CommitBuildingUpgradeParams{
		BuildingID:   nextLevel.NextBuildingID,
		CurrentLevel: ownedBuilding.CurrentLevel + 1,
		ID:           pgPlacementID,
	}
	if err := queries.CommitBuildingUpgrade(ctx, upgradeArgs); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to apply upgrade to building")
	}

	return map[string]interface{}{
		"message":       "Building upgraded successfully!",
		"new_level":     ownedBuilding.CurrentLevel + 1,
		"building_name": nextLevel.Name,
	}, http.StatusOK, nil
}

func GetBuildingUpgradeCost(ctx context.Context, queries *db.Queries, pgPlayerID pgtype.UUID, placementID string) (map[string]interface{}, int, error) {
	parsedPlacementUUID, err := uuid.Parse(placementID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid placement ID format")
	}
	pgPlacementID := pgtype.UUID{Bytes: parsedPlacementUUID, Valid: true}
	_ = queries.AutoCompleteBuildings(ctx, pgPlayerID)

	ownedArgs := db.GetOwnedBuildingDetailsParams{
		ID:       pgPlacementID,
		PlayerID: pgPlayerID,
	}
	ownedBuilding, err := queries.GetOwnedBuildingDetails(ctx, ownedArgs)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Building not found or you don't own it")
	}

	if !ownedBuilding.IsBuilt.Bool {
		return nil, http.StatusBadRequest, fmt.Errorf("Building is currently under construction or upgrade")
	}

	nextLevelArgs := db.GetNextLevelBuildingParams{
		BuildingType: ownedBuilding.BuildingType,
		Level:        ownedBuilding.CurrentLevel + 1,
	}
	nextLevel, err := queries.GetNextLevelBuilding(ctx, nextLevelArgs)
	if err != nil {
		return map[string]interface{}{
			"is_max_level": true,
			"message":      "Building is already at max level!",
		}, http.StatusOK, nil
	}

	return map[string]interface{}{
		"is_max_level": false,
		"cost_type":    strings.ToLower(string(nextLevel.CostType)),
		"cost":         nextLevel.BuildCost,
		"next_level":   ownedBuilding.CurrentLevel + 1,
		"name":         nextLevel.Name,
	}, http.StatusOK, nil
}

func UpgradePlayerVillage(ctx context.Context, pool *pgxpool.Pool, pgPlayerID pgtype.UUID) (map[string]interface{}, int, error) {
	queries := db.New(pool)

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Database error")
	}
	defer tx.Rollback(ctx)
	qtx := queries.WithTx(tx)

	player, err := qtx.GetPlayerProfile(ctx, pgPlayerID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Player profile not found")
	}

	goldCost := player.VillageLevel * 10000

	if player.VillageLevel >= 4 {
		return nil, http.StatusBadRequest, fmt.Errorf("Your Village has already reached maximum level!")
	}

	if player.GoldCoins < goldCost {
		return nil, http.StatusPaymentRequired, fmt.Errorf("Insufficient Gold Coins to upgrade your Village!")
	}

	err = qtx.UpgradeVillageLevel(ctx, db.UpgradeVillageLevelParams{
		GoldCost: goldCost,
		PlayerID: pgPlayerID,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to upgrade village level")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to finalize village upgrade")
	}

	return map[string]interface{}{
		"message":           "Congratulations! Your Village level has increased!",
		"old_village_level": player.VillageLevel,
		"new_village_level": player.VillageLevel + 1,
		"gold_spent":        goldCost,
		"remaining_gold":    player.GoldCoins - goldCost,
	}, http.StatusOK, nil
}

func GetVillageUpgradeCost(ctx context.Context, pool *pgxpool.Pool, pgPlayerID pgtype.UUID) (map[string]interface{}, int, error) {
	queries := db.New(pool)
	player, err := queries.GetPlayerProfile(ctx, pgPlayerID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Player profile not found")
	}

	goldCost := player.VillageLevel * 10000

	if player.VillageLevel >= 4 {
		return map[string]interface{}{
			"is_max_level": true,
			"message":      "Your Village has already reached maximum level!",
		}, http.StatusOK, nil
	}

	return map[string]interface{}{
		"is_max_level": false,
		"cost_type":    "gold",
		"cost":         goldCost,
		"next_level":   player.VillageLevel + 1,
	}, http.StatusOK, nil
}

// CalculateVillageUpgradeCost checks maximum level restrictions and computes the gold cost for upgrading the Town Hall.
func CalculateVillageUpgradeCost(currentLevel int32) (isMax bool, nextLevel int32, goldCost int32) {
	if currentLevel >= 4 {
		return true, currentLevel, 0
	}
	return false, currentLevel + 1, currentLevel * 10000
}

// ValidateUpgradeEligibility checks whether a player satisfies level and resource requirements for building progression.
func ValidateUpgradeEligibility(villageLevel int32, requiredLevel int32, playerGold int32, playerElixir int32, costType string, buildCost int32) (bool, string) {
	if villageLevel < requiredLevel {
		return false, "Your Village Level is too low to upgrade this building"
	}
	if strings.ToLower(costType) == "gold" && playerGold < buildCost {
		return false, "Insufficient resources to upgrade"
	}
	if strings.ToLower(costType) == "elixir" && playerElixir < buildCost {
		return false, "Insufficient resources to upgrade"
	}
	return true, "OK"
}
