package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func CollectResources(ctx context.Context, queries *db.Queries, pgPlayerID pgtype.UUID) (map[string]interface{}, int, error) {
	buildings, err := queries.GetResourceBuildingsForCollection(ctx, pgPlayerID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to fetch resource buildings")
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
		return map[string]interface{}{
			"message": "No new resources to collect right now. Check back later!",
		}, http.StatusOK, nil
	}

	profile, err := queries.GetPlayerProfile(ctx, pgPlayerID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to fetch player profile")
	}

	storages, err := queries.GetTotalStorageCapacity(ctx, pgPlayerID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to fetch storage capacity")
	}

	var goldCap int32 = 1000
	var elixirCap int32 = 1000

	for _, s := range storages {
		if s.ResourceType == db.ResourceTypeGold {
			goldCap += s.TotalCapacity
		} else if s.ResourceType == db.ResourceTypeElixir {
			elixirCap += s.TotalCapacity
		}
	}

	goldSpace := goldCap - profile.GoldCoins
	elixirSpace := elixirCap - profile.Elixir

	if goldSpace < 0 {
		goldSpace = 0
	}
	if elixirSpace < 0 {
		elixirSpace = 0
	}

	if totalGoldGenerated > goldSpace {
		totalGoldGenerated = goldSpace
	}
	if totalElixirGenerated > elixirSpace {
		totalElixirGenerated = elixirSpace
	}

	if totalGoldGenerated == 0 && totalElixirGenerated == 0 {
		return map[string]interface{}{
			"message": "Storage is full! Upgrade your Gold/Elixir Storage to collect more.",
		}, http.StatusOK, nil
	}

	addParams := db.AddPlayerResourcesParams{
		GoldCoins: totalGoldGenerated,
		Elixir:    totalElixirGenerated,
		ID:        pgPlayerID,
	}

	updatedProfile, err := queries.AddPlayerResources(ctx, addParams)
	if err != nil {
		fmt.Printf("Error adding resources: %v\n", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to add resources to profile")
	}

	err = queries.ResetCollectionTime(ctx, pgPlayerID)
	if err != nil {
		fmt.Printf("Error resetting timers: %v\n", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to reset collection timers")
	}

	return map[string]interface{}{
		"message":       "Resources collected successfully!",
		"gold_looted":   totalGoldGenerated,
		"elixir_looted": totalElixirGenerated,
		"new_balance": map[string]int32{
			"gold":   updatedProfile.GoldCoins,
			"elixir": updatedProfile.Elixir,
		},
		"storage_caps": map[string]interface{}{
			"gold_cap":   goldCap,
			"elixir_cap": elixirCap,
		},
	}, http.StatusOK, nil
}
