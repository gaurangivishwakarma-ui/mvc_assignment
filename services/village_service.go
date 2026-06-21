package services

import (
	"context"
	"fmt"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetVillage(ctx context.Context, queries *db.Queries, pgPlayerID pgtype.UUID) (map[string]interface{}, int, error) {
	profile, err := queries.GetPlayerProfile(ctx, pgPlayerID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to fetch player profile")
	}

	buildings, err := queries.GetPlayerBuildings(ctx, pgPlayerID)
	if err != nil {
		buildings = []db.BuildingsOwned{}
	}

	return map[string]interface{}{
		"profile":   profile,
		"buildings": buildings,
	}, http.StatusOK, nil
}
