package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/gaurangi/mvc_assignment/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCalculateCombatPower_GiantsAndGoblins(t *testing.T) {
	tests := []struct {
		name             string
		initialAttacker  int32
		initialDefense   int32
		deployedTroops   []models.DeployedTroop
		expectedAttacker int32
		expectedDefense  int32
	}{
		{
			name:             "Standard troops with no special ability modification",
			initialAttacker:  500,
			initialDefense:   1000,
			deployedTroops:   []models.DeployedTroop{{TroopType: "barbarian", Quantity: 10, Level: 1}},
			expectedAttacker: 500,
			expectedDefense:  1000,
		},
		{
			name:             "Giants reduce defense power significantly",
			initialAttacker:  600,
			initialDefense:   1500,
			deployedTroops:   []models.DeployedTroop{{TroopType: "giant", Quantity: 5, Level: 1}}, // reduction = 5*1*200 = 1000
			expectedAttacker: 600,
			expectedDefense:  500,
		},
		{
			name:             "Giants cap defense power reduction at 1 to prevent zero division",
			initialAttacker:  500,
			initialDefense:   300,
			deployedTroops:   []models.DeployedTroop{{TroopType: "giant", Quantity: 4, Level: 1}}, // reduction = 800
			expectedAttacker: 500,
			expectedDefense:  1,
		},
		{
			name:             "Goblins increase attack power for resource raids",
			initialAttacker:  400,
			initialDefense:   800,
			deployedTroops:   []models.DeployedTroop{{TroopType: "goblin", Quantity: 10, Level: 2}}, // bonus = 10*2*85 = 1700
			expectedAttacker: 2100,
			expectedDefense:  800,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			att, def := CalculateCombatPower(tt.initialAttacker, tt.initialDefense, tt.deployedTroops)
			if att != tt.expectedAttacker {
				t.Errorf("expected attacker power %d, got %d", tt.expectedAttacker, att)
			}
			if def != tt.expectedDefense {
				t.Errorf("expected defense power %d, got %d", tt.expectedDefense, def)
			}
		})
	}
}

func TestCalculateCappedLoot(t *testing.T) {
	tests := []struct {
		name             string
		stolenAmount     int32
		storageCap       int32
		currentBalance   int32
		expectedLoot     int32
		expectedIsCapped bool
	}{
		{
			name:             "Loot within available storage limits",
			stolenAmount:     500,
			storageCap:       10000,
			currentBalance:   2000,
			expectedLoot:     500,
			expectedIsCapped: false,
		},
		{
			name:             "Loot exceeds available storage capacity",
			stolenAmount:     8000,
			storageCap:       10000,
			currentBalance:   7000, // available space = 3000
			expectedLoot:     3000,
			expectedIsCapped: true,
		},
		{
			name:             "Storage is completely full before raid",
			stolenAmount:     1500,
			storageCap:       5000,
			currentBalance:   5500, // over capacity
			expectedLoot:     0,
			expectedIsCapped: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualLoot, capped := CalculateCappedLoot(tt.stolenAmount, tt.storageCap, tt.currentBalance)
			if actualLoot != tt.expectedLoot || capped != tt.expectedIsCapped {
				t.Errorf("expected (%d, %t), got (%d, %t)", tt.expectedLoot, tt.expectedIsCapped, actualLoot, capped)
			}
		})
	}
}

func TestAttackOpponent_InvalidUUID(t *testing.T) {
	ctx := context.Background()
	req := models.InstantBattleRequest{
		OpponentID: "not-a-valid-uuid",
	}
	_, status, err := AttackOpponent(ctx, nil, pgtype.UUID{Valid: true}, req)
	if status != http.StatusBadRequest || err == nil {
		t.Errorf("expected HTTP 400 Bad Request for invalid opponent UUID, got %d", status)
	}
}
