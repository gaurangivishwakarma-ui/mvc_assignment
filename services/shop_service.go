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

func CheckOverlap(x1, y1, w1, h1 int32, x2, y2, w2, h2 int32) bool {
	return !(x1+w1 <= x2 || x1 >= x2+w2 || y1+h1 <= y2 || y1 >= y2+h2)
}

func PurchaseBuilding(ctx context.Context, queries *db.Queries, pgPlayerID pgtype.UUID, req models.PurchaseBuildingRequest) (map[string]interface{}, int, error) {
	catalogBuilding, err := queries.GetBuildingByID(ctx, req.BuildingID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Building not found in catalog")
	}

	existingBuildings, err := queries.GetPlayerBuildingsWithDimensions(ctx, pgPlayerID)
	if err == nil {
		for _, b := range existingBuildings {
			overlaps := CheckOverlap(
				req.XCoords, req.YCoords, catalogBuilding.Width, catalogBuilding.Breadth,
				b.XCoords, b.YCoords, b.Width, b.Breadth,
			)

			if overlaps {
				return nil, http.StatusConflict, fmt.Errorf("Cannot place building here. Space is occupied!")
			}
		}
	}

	deductParams := db.DeductResourcesParams{
		GoldCoins: catalogBuilding.Cost,
		Elixir:    0,
		ID:        pgPlayerID,
	}

	_, err = queries.DeductResources(ctx, deductParams)
	if err != nil {
		return nil, http.StatusPaymentRequired, fmt.Errorf("Insufficient resources to purchase this building")
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

	newBuilding, err := queries.PlaceBuilding(ctx, placeParams)
	if err != nil {
		fmt.Printf("DB Error placing building: %v\n", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to place building")
	}

	return map[string]interface{}{
		"message":  "Building purchased and placed successfully!",
		"building": newBuilding,
	}, http.StatusCreated, nil
}

func CompleteBuild(ctx context.Context, queries *db.Queries, pgPlayerID pgtype.UUID, req models.CompleteBuildRequest) (map[string]interface{}, int, error) {
	parsedUUID, err := uuid.Parse(req.PlacementID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid placement ID format")
	}
	pgBuildingID := pgtype.UUID{Bytes: parsedUUID, Valid: true}

	rowsAffected, err := queries.MarkBuildingAsBuilt(ctx, db.MarkBuildingAsBuiltParams{
		ID:       pgBuildingID,
		PlayerID: pgPlayerID,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to complete building construction")
	}

	if rowsAffected == 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("Building is not ready yet (5 seconds must pass after purchase)")
	}

	return map[string]interface{}{
		"message":      "Building construction complete!",
		"placement_id": req.PlacementID,
	}, http.StatusOK, nil
}
