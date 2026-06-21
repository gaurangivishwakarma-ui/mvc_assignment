package services

import (
	"context"
	"fmt"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetDashboard(ctx context.Context, queries *db.Queries, pgPlayerID pgtype.UUID) (models.DashboardResponse, int, error) {
	rows, err := queries.GetDashboardData(ctx, pgPlayerID)
	if err != nil || len(rows) == 0 {
		return models.DashboardResponse{}, http.StatusInternalServerError, fmt.Errorf("Profile data unreachable")
	}

	firstRow := rows[0]
	response := models.DashboardResponse{
		Username:     firstRow.Username,
		VillageLevel: firstRow.VillageLevel,
		XPPoints:     firstRow.XpPoints,
		Balances: map[string]int32{
			"gold":   firstRow.GoldCoins,
			"elixir": firstRow.Elixir,
		},
		Buildings: []models.BuildingInventory{},
	}

	for _, row := range rows {
		if row.PlacementID.Valid {
			var uidStr string
			bytes := row.PlacementID.Bytes
			if parsed, err := uuid.FromBytes(bytes[:]); err == nil {
				uidStr = parsed.String()
			}

			response.Buildings = append(response.Buildings, models.BuildingInventory{
				PlacementID:  uidStr,
				BuildingID:   row.BuildingID.Int32,
				BuildingType: row.BuildingType.String,
				CurrentLevel: row.CurrentLevel.Int32,
				XCoords:      row.XCoords.Int32,
				YCoords:      row.YCoords.Int32,
				IsBuilt:      row.IsBuilt.Bool,
			})
		}
	}

	return response, http.StatusOK, nil
}
