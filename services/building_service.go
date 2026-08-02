package services

import (
	"context"
	"fmt"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MoveBuilding(ctx context.Context, pool *pgxpool.Pool, pgPlayerID pgtype.UUID, req models.MoveBuildingRequest) (map[string]interface{}, int, error) {
	parsedBuildingUUID, err := uuid.Parse(req.OwnedBuildingID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid building ID format")
	}
	pgBuildingID := pgtype.UUID{Bytes: parsedBuildingUUID, Valid: true}

	queries := db.New(pool)
	_ = queries.AutoCompleteBuildings(ctx, pgPlayerID)

	info, err := queries.GetOwnedBuildingPositionInfo(ctx, pgBuildingID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Owned building not found")
	}

	if info.PlayerID != pgPlayerID {
		return nil, http.StatusForbidden, fmt.Errorf("Unauthorized: You do not own this structure")
	}

	if !info.IsBuilt.Bool {
		return nil, http.StatusBadRequest, fmt.Errorf("Cannot move building: 5 seconds must pass for creating or upgrading it")
	}

	const maxMapGrid int32 = 40
	if req.NewX < 0 || req.NewY < 0 || (req.NewX+info.Width) > maxMapGrid || (req.NewY+info.Breadth) > maxMapGrid {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid coordinates! Structure exceeds village map boundaries.")
	}

	allBuildings, err := queries.GetPlayerBuildingsWithDimensions(ctx, pgPlayerID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to fetch village layout")
	}

	isOverlapping := false
	for _, b := range allBuildings {
		if b.ID == pgBuildingID {
			continue
		}
		if req.NewX < (b.XCoords+b.Width) &&
			(req.NewX+info.Width) > b.XCoords &&
			req.NewY < (b.YCoords+b.Breadth) &&
			(req.NewY+info.Breadth) > b.YCoords {
			isOverlapping = true
			break
		}
	}

	if isOverlapping {
		return nil, http.StatusConflict, fmt.Errorf("Invalid placement! Building overlaps with another structure.")
	}

	err = queries.UpdateBuildingPosition(ctx, db.UpdateBuildingPositionParams{
		ID:       pgBuildingID,
		XCoords:  req.NewX,
		YCoords:  req.NewY,
		PlayerID: pgPlayerID,
	})
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to update layout positions")
	}

	return map[string]interface{}{
		"status": "Building moved successfully",
		"new_x":  req.NewX,
		"new_y":  req.NewY,
	}, http.StatusOK, nil
}
