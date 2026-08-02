package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/gaurangi/mvc_assignment/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestValidateTroopTraining(t *testing.T) {
	tests := []struct {
		name              string
		currentUsed       int32
		maxCapacity       int32
		housingPerTroop   int32
		requestedQuantity int32
		expectedAllowed   bool
		expectedFitCount  int32
	}{
		{
			name:              "Training barbarians within available housing space",
			currentUsed:       10,
			maxCapacity:       30,
			housingPerTroop:   1,
			requestedQuantity: 15,
			expectedAllowed:   true,
			expectedFitCount:  15,
		},
		{
			name:              "Training giants exceeding available camp space",
			currentUsed:       20,
			maxCapacity:       30, // 10 space left
			housingPerTroop:   5,  // Giants take 5 housing space each
			requestedQuantity: 4,  // requires 20 space
			expectedAllowed:   false,
			expectedFitCount:  2, // only 2 giants can fit in 10 remaining space
		},
		{
			name:              "Army camp already completely full",
			currentUsed:       50,
			maxCapacity:       50,
			housingPerTroop:   1,
			requestedQuantity: 1,
			expectedAllowed:   false,
			expectedFitCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, fitCount := ValidateTroopTraining(tt.currentUsed, tt.maxCapacity, tt.housingPerTroop, tt.requestedQuantity)
			if allowed != tt.expectedAllowed || fitCount != tt.expectedFitCount {
				t.Errorf("expected (%t, %d fit), got (%t, %d fit)", tt.expectedAllowed, tt.expectedFitCount, allowed, fitCount)
			}
		})
	}
}

func TestTrainTroops_InvalidPayload(t *testing.T) {
	ctx := context.Background()
	req := models.TrainTroopRequest{
		TroopType: "barbarian",
		Quantity:  0, // Invalid quantity
		Level:     1,
	}
	_, status, err := TrainTroops(ctx, nil, pgtype.UUID{Valid: true}, req)
	if status != http.StatusBadRequest || err == nil {
		t.Errorf("expected HTTP 400 for 0 quantity training request, got %d", status)
	}
}
