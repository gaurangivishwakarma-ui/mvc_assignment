package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PurchaseBuildingRequest struct {
	BuildingID int32 `json:"building_id"`
	XCoords    int32 `json:"x_coords"`
	YCoords    int32 `json:"y_coords"`
}

func checkOverlap(x1, y1, w1, h1 int32, x2, y2, w2, h2 int32) bool {
	return !(x1+w1 <= x2 || x1 >= x2+w2 || y1+h1 <= y2 || y1 >= y2+h2)
}

func PurchaseBuilding(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req PurchaseBuildingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON request", http.StatusBadRequest)
			return
		}

		pgPlayerID := r.Context().Value(middleware.PlayerIDKey).(pgtype.UUID)

		catalogBuilding, err := queries.GetBuildingByID(r.Context(), req.BuildingID)
		if err != nil {
			http.Error(w, "Building not found in catalog", http.StatusNotFound)
			return
		}

		existingBuildings, err := queries.GetPlayerBuildingsWithDimensions(r.Context(), pgPlayerID)
		if err == nil {
			for _, b := range existingBuildings {
				overlaps := checkOverlap(
					req.XCoords, req.YCoords, catalogBuilding.Width, catalogBuilding.Breadth,
					b.XCoords, b.YCoords, b.Width, b.Breadth,
				)

				if overlaps {
					http.Error(w, "Cannot place building here. Space is occupied!", http.StatusConflict)
					return
				}
			}
		}

		deductParams := db.DeductResourcesParams{
			GoldCoins: catalogBuilding.Cost,
			Elixir:    0,
			ID:        pgPlayerID,
		}

		_, err = queries.DeductResources(r.Context(), deductParams)
		if err != nil {
			http.Error(w, "Insufficient resources to purchase this building", http.StatusPaymentRequired)
			return
		}

		newBuildingID := uuid.New()
		pgBuildingID := pgtype.UUID{Bytes: newBuildingID, Valid: true}

		placeParams := db.PlaceBuildingParams{
			ID:           pgBuildingID,
			PlayerID:     pgPlayerID,
			BuildingType: catalogBuilding.BuildingType,
			BuildingID:   catalogBuilding.ID,
			XCoords:      req.XCoords,
			YCoords:      req.YCoords,
		}

		newBuilding, err := queries.PlaceBuilding(r.Context(), placeParams)

		if err != nil {
			fmt.Printf("DB Error placing building: %v\n", err)
			http.Error(w, "Failed to place building", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Building purchased and placed successfully!",
			"building": newBuilding,
		})
	}
}
