package services

import (
	"testing"
)

func TestCalculateResourceGeneration(t *testing.T) {
	tests := []struct {
		name              string
		hoursPassed       float64
		hourlyRate        float64
		currentBalance    int32
		storageCap        int32
		expectedGenerated int32
	}{
		{
			name:              "Normal resource gathering below capacity",
			hoursPassed:       5.0,
			hourlyRate:        100.0,
			currentBalance:    2000,
			storageCap:        10000,
			expectedGenerated: 500, // 5 * 100 = 500
		},
		{
			name:              "Collection capped when reaching max storage capacity",
			hoursPassed:       10.0,
			hourlyRate:        200.0, // generates 2000
			currentBalance:    9000,
			storageCap:        10000,
			expectedGenerated: 1000, // limited to 1000 remaining space
		},
		{
			name:              "Zero gathering when storage is completely full",
			hoursPassed:       12.0,
			hourlyRate:        300.0,
			currentBalance:    15000,
			storageCap:        10000, // balance exceeds cap
			expectedGenerated: 0,
		},
		{
			name:              "Zero gathering for very short duration",
			hoursPassed:       0.0,
			hourlyRate:        500.0,
			currentBalance:    1000,
			storageCap:        10000,
			expectedGenerated: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CalculateResourceGeneration(tt.hoursPassed, tt.hourlyRate, tt.currentBalance, tt.storageCap)
			if actual != tt.expectedGenerated {
				t.Errorf("expected %d generated resources, got %d", tt.expectedGenerated, actual)
			}
		})
	}
}
