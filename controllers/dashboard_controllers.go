package controllers

import (
	"encoding/json"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type BuildingInventory struct {
	PlacementID  string `json:"placement_id,omitempty"`
	BuildingID   int32  `json:"building_id,omitempty"`
	BuildingType string `json:"building_type,omitempty"`
	CurrentLevel int32  `json:"current_level,omitempty"`
	XCoords      int32  `json:"x_coords"`
	YCoords      int32  `json:"y_coords"`
	IsBuilt      bool   `json:"is_built"`
}

type DashboardResponse struct {
	Username     string              `json:"username"`
	VillageLevel int32               `json:"village_level"`
	XPPoints     int32               `json:"xp_points"`
	Balances     map[string]int32    `json:"balances"`
	Buildings    []BuildingInventory `json:"buildings_placed"`
}

func GetDashboard(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		playerIDStr := r.Context().Value(middleware.PlayerIDKey).(string)
		parsedUUID, _ := uuid.Parse(playerIDStr)
		pgPlayerID := pgtype.UUID{Bytes: parsedUUID, Valid: true}

		rows, err := queries.GetDashboardData(r.Context(), pgPlayerID)
		if err != nil || len(rows) == 0 {
			http.Error(w, "Profile data unreachable", http.StatusInternalServerError)
			return
		}

		firstRow := rows[0]
		response := DashboardResponse{
			Username:     firstRow.Username,
			VillageLevel: firstRow.VillageLevel,
			XPPoints:     firstRow.XpPoints,
			Balances: map[string]int32{
				"gold":   firstRow.GoldCoins,
				"elixir": firstRow.Elixir,
			},
			Buildings: []BuildingInventory{},
		}

		for _, row := range rows {
			if row.PlacementID.Valid {
				var uidStr string
				bytes := row.PlacementID.Bytes
				if parsed, err := uuid.FromBytes(bytes[:]); err == nil {
					uidStr = parsed.String()
				}

				response.Buildings = append(response.Buildings, BuildingInventory{
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
