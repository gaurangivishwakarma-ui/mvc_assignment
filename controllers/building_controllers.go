package controllers

import (
	"encoding/json"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MoveBuildingRequest struct {
	OwnedBuildingID string `json:"owned_building_id"`
	NewX            int32  `json:"new_x"`
	NewY            int32  `json:"new_y"`
}

func MoveBuilding(pool *pgxpool.Pool) http.HandlerFunc {
	queries := db.New(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		playerIDStr := r.Context().Value(middleware.PlayerIDKey).(string)
		parsedPlayerUUID, _ := uuid.Parse(playerIDStr)
		pgPlayerID := pgtype.UUID{Bytes: parsedPlayerUUID, Valid: true}

		var req MoveBuildingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		parsedBuildingUUID, err := uuid.Parse(req.OwnedBuildingID)
		if err != nil {
			http.Error(w, "Invalid building ID format", http.StatusBadRequest)
			return
		}
		pgBuildingID := pgtype.UUID{Bytes: parsedBuildingUUID, Valid: true}

		info, err := queries.GetOwnedBuildingPositionInfo(r.Context(), pgBuildingID)
		if err != nil {
			http.Error(w, "Owned building not found", http.StatusNotFound)
			return
		}

		if info.PlayerID != pgPlayerID {
			http.Error(w, "Unauthorized: You do not own this structure", http.StatusForbidden)
			return
		}

		const maxMapGrid int32 = 40
		if req.NewX < 0 || req.NewY < 0 || (req.NewX+info.Width) > maxMapGrid || (req.NewY+info.Breadth) > maxMapGrid {
			http.Error(w, "Invalid coordinates! Structure exceeds village map boundaries.", http.StatusBadRequest)
			return
		}

		allBuildings, err := queries.GetPlayerBuildingsWithDimensions(r.Context(), pgPlayerID)
		if err != nil {
			http.Error(w, "Failed to fetch village layout", http.StatusInternalServerError)
			return
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
			http.Error(w, "Invalid placement! Building overlaps with another structure.", http.StatusConflict)
			return
		}

		err = queries.UpdateBuildingPosition(r.Context(), db.UpdateBuildingPositionParams{
			ID:       pgBuildingID,
			XCoords:  req.NewX,
			YCoords:  req.NewY,
			PlayerID: pgPlayerID,
		})
		if err != nil {
			http.Error(w, "Failed to update layout positions", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "Building moved successfully",
			"new_x":  req.NewX,
			"new_y":  req.NewY,
		})
	}
}
