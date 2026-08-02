package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/gaurangi/mvc_assignment/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCheckOverlap(t *testing.T) {
	tests := []struct {
		name     string
		x1, y1   int32
		w1, h1   int32
		x2, y2   int32
		w2, h2   int32
		expected bool
	}{
		{
			name: "No overlap - buildings side by side horizontally",
			x1:   0, y1: 0, w1: 3, h1: 3,
			x2: 3, y2: 0, w2: 3, h2: 3,
			expected: false,
		},
		{
			name: "No overlap - buildings separated vertically",
			x1:   5, y1: 5, w1: 4, h1: 4,
			x2: 5, y2: 10, w2: 3, h2: 3,
			expected: false,
		},
		{
			name: "Direct overlap - identical coordinates",
			x1:   10, y1: 10, w1: 3, h1: 3,
			x2: 10, y2: 10, w2: 3, h2: 3,
			expected: true,
		},
		{
			name: "Partial intersection horizontally and vertically",
			x1:   10, y1: 10, w1: 4, h1: 4, // spans x:10..14, y:10..14
			x2: 12, y2: 12, w2: 3, h2: 3, // spans x:12..15, y:12..15
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CheckOverlap(tt.x1, tt.y1, tt.w1, tt.h1, tt.x2, tt.y2, tt.w2, tt.h2)
			if actual != tt.expected {
				t.Errorf("expected overlap %t, got %t", tt.expected, actual)
			}
		})
	}
}

func TestMoveBuilding_InvalidUUID(t *testing.T) {
	ctx := context.Background()
	req := models.MoveBuildingRequest{
		OwnedBuildingID: "invalid-uuid",
		NewX:            5,
		NewY:            5,
	}
	_, status, err := MoveBuilding(ctx, nil, pgtype.UUID{Valid: true}, req)
	if status != http.StatusBadRequest || err == nil {
		t.Errorf("expected HTTP 400 for malformed owned building UUID, got %d", status)
	}
}
