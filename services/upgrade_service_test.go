package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/gaurangi/mvc_assignment/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCalculateVillageUpgradeCost(t *testing.T) {
	tests := []struct {
		name         string
		currentLevel int32
		expectedMax  bool
		expectedNext int32
		expectedCost int32
	}{
		{"Upgrade Lv 1 to Lv 2", 1, false, 2, 10000},
		{"Upgrade Lv 2 to Lv 3", 2, false, 3, 20000},
		{"Upgrade Lv 3 to Lv 4", 3, false, 4, 30000},
		{"Max Level reached (Lv 4)", 4, true, 4, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isMax, nextLvl, cost := CalculateVillageUpgradeCost(tt.currentLevel)
			if isMax != tt.expectedMax || nextLvl != tt.expectedNext || cost != tt.expectedCost {
				t.Errorf("expected (%t, %d, %d), got (%t, %d, %d)", tt.expectedMax, tt.expectedNext, tt.expectedCost, isMax, nextLvl, cost)
			}
		})
	}
}

func TestValidateUpgradeEligibility(t *testing.T) {
	tests := []struct {
		name          string
		villageLevel  int32
		requiredLevel int32
		playerGold    int32
		playerElixir  int32
		costType      string
		buildCost     int32
		expectedOk    bool
	}{
		{"Valid upgrade with sufficient gold", 2, 1, 5000, 1000, "gold", 2000, true},
		{"Valid upgrade with sufficient elixir", 2, 2, 1000, 4000, "elixir", 3000, true},
		{"Blocked by insufficient village level", 1, 2, 10000, 10000, "gold", 500, false},
		{"Blocked by insufficient gold", 2, 1, 100, 10000, "gold", 500, false},
		{"Blocked by insufficient elixir", 2, 1, 10000, 200, "elixir", 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, _ := ValidateUpgradeEligibility(tt.villageLevel, tt.requiredLevel, tt.playerGold, tt.playerElixir, tt.costType, tt.buildCost)
			if ok != tt.expectedOk {
				t.Errorf("expected eligibility %t, got %t", tt.expectedOk, ok)
			}
		})
	}
}

func TestUpgradeBuilding_InvalidUUID(t *testing.T) {
	ctx := context.Background()
	req := models.UpgradeRequest{
		PlacementID: "invalid-placement-uuid",
	}
	_, status, err := UpgradeBuilding(ctx, nil, pgtype.UUID{Valid: true}, req)
	if status != http.StatusBadRequest || err == nil {
		t.Errorf("expected HTTP 400 Bad Request for invalid placement ID, got %d", status)
	}
}
