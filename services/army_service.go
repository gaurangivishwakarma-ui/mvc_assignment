package services

import (
	"context"
	"fmt"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TrainTroops(ctx context.Context, pool *pgxpool.Pool, pgPlayerID pgtype.UUID, req models.TrainTroopRequest) (map[string]interface{}, int, error) {
	queries := db.New(pool)

	if req.Quantity <= 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid request payload or quantity")
	}

	if req.Level <= 0 {
		req.Level = 1
	}

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

	troopArgs := db.GetTroopDetailsParams{
		TroopType: req.TroopType,
		LevelReq:  req.Level,
	}
	troopInfo, err := qtx.GetTroopDetails(ctx, troopArgs)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Troop type or level not found in catalog")
	}

	if player.VillageLevel < troopInfo.LevelReq {
		return nil, http.StatusForbidden, fmt.Errorf("Village level is too low to train this troop!")
	}

	maxCapacity, _ := qtx.GetTotalHousingCapacity(ctx, pgPlayerID)
	usedCapacity, _ := qtx.GetCurrentHousingUsed(ctx, pgPlayerID)

	spaceNeeded := troopInfo.HousingSpace * req.Quantity
	if usedCapacity+spaceNeeded > maxCapacity {
		return nil, http.StatusForbidden, fmt.Errorf("Not enough Army Camp housing space!")
	}

	totalCost := troopInfo.ElixirCost * req.Quantity
	payArgs := db.PayForTroopTrainingParams{
		Elixir: totalCost,
		ID:     pgPlayerID,
	}
	if err := qtx.PayForTroopTraining(ctx, payArgs); err != nil {
		return nil, http.StatusPaymentRequired, fmt.Errorf("Insufficient Elixir to train these troops")
	}

	addArgs := db.AddTroopsToArmyParams{
		PlayerID:     pgPlayerID,
		TroopType:    req.TroopType,
		CurrentLevel: req.Level,
		Quantity:     req.Quantity,
	}
	if err := qtx.AddTroopsToArmy(ctx, addArgs); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to add troops to active army")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to finalize training")
	}

	return map[string]interface{}{
		"message":    "Troops instantly trained and ready for battle!",
		"troop_type": req.TroopType,
		"level":      req.Level,
		"quantity":   req.Quantity,
		"cost_paid":  totalCost,
		"space_used": usedCapacity + spaceNeeded,
		"space_max":  maxCapacity,
	}, http.StatusOK, nil
}

func GetArmyCatalog(ctx context.Context, queries *db.Queries) (map[string]interface{}, int, error) {
	catalog, err := queries.GetTroopCatalog(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to fetch troop catalog")
	}
	return map[string]interface{}{
		"troops_available": catalog,
	}, http.StatusOK, nil
}

func GetArmyStatus(ctx context.Context, queries *db.Queries, pgPlayerID pgtype.UUID) (map[string]interface{}, int, error) {
	rows, err := queries.GetArmyStatus(ctx, pgPlayerID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to fetch army status")
	}

	var army []models.ArmyStatusResponse
	for _, row := range rows {
		army = append(army, models.ArmyStatusResponse{
			TroopType:    row.TroopType,
			CurrentLevel: row.CurrentLevel,
			Quantity:     row.Quantity,
		})
	}

	if army == nil {
		army = []models.ArmyStatusResponse{}
	}

	return map[string]interface{}{
		"army": army,
	}, http.StatusOK, nil
}

// ValidateTroopTraining checks whether an army camp has adequate remaining housing space for a new troop batch.
func ValidateTroopTraining(currentUsedSpace int32, maxCampCapacity int32, housingPerTroop int32, quantity int32) (bool, int32) {
	requiredSpace := housingPerTroop * quantity
	availableSpace := maxCampCapacity - currentUsedSpace
	if availableSpace < 0 {
		availableSpace = 0
	}
	if requiredSpace > availableSpace {
		if housingPerTroop <= 0 {
			return false, 0
		}
		return false, availableSpace / housingPerTroop
	}
	return true, quantity
}
