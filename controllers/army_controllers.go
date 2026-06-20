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

type TrainTroopRequest struct {
	TroopType string `json:"troop_type"`
	Level     int32  `json:"level"`
	Quantity  int32  `json:"quantity"`
}

func TrainTroops(pool *pgxpool.Pool) http.HandlerFunc {
	queries := db.New(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		playerIDStr := r.Context().Value(middleware.PlayerIDKey).(string)
		parsedPlayerUUID, _ := uuid.Parse(playerIDStr)
		pgPlayerID := pgtype.UUID{Bytes: parsedPlayerUUID, Valid: true}

		var req TrainTroopRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Quantity <= 0 {
			http.Error(w, "Invalid request payload or quantity", http.StatusBadRequest)
			return
		}

		if req.Level <= 0 {
			req.Level = 1
		}

		tx, err := pool.Begin(r.Context())
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context())
		qtx := queries.WithTx(tx)

		player, err := qtx.GetPlayerProfile(r.Context(), pgPlayerID)
		if err != nil {
			http.Error(w, "Player profile not found", http.StatusNotFound)
			return
		}

		troopArgs := db.GetTroopDetailsParams{
			TroopType: req.TroopType,
			LevelReq:  req.Level,
		}
		troopInfo, err := qtx.GetTroopDetails(r.Context(), troopArgs)
		if err != nil {
			http.Error(w, "Troop type or level not found in catalog", http.StatusNotFound)
			return
		}

		if player.VillageLevel < troopInfo.LevelReq {
			http.Error(w, "Village level is too low to train this troop!", http.StatusForbidden)
			return
		}

		maxCapacity, _ := qtx.GetTotalHousingCapacity(r.Context(), pgPlayerID)
		usedCapacity, _ := qtx.GetCurrentHousingUsed(r.Context(), pgPlayerID)

		spaceNeeded := troopInfo.HousingSpace * req.Quantity
		if usedCapacity+spaceNeeded > maxCapacity {
			http.Error(w, "Not enough Army Camp housing space!", http.StatusForbidden)
			return
		}

		totalCost := troopInfo.ElixirCost * req.Quantity
		payArgs := db.PayForTroopTrainingParams{
			Elixir: totalCost,
			ID:     pgPlayerID,
		}
		if err := qtx.PayForTroopTraining(r.Context(), payArgs); err != nil {
			http.Error(w, "Insufficient Elixir to train these troops", http.StatusPaymentRequired)
			return
		}

		addArgs := db.AddTroopsToArmyParams{
			PlayerID:     pgPlayerID,
			TroopType:    req.TroopType,
			CurrentLevel: req.Level,
			Quantity:     req.Quantity,
		}
		if err := qtx.AddTroopsToArmy(r.Context(), addArgs); err != nil {
			http.Error(w, "Failed to add troops to active army", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(r.Context()); err != nil {
			http.Error(w, "Failed to finalize training", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":    "Troops instantly trained and ready for battle!",
			"troop_type": req.TroopType,
			"level":      req.Level,
			"quantity":   req.Quantity,
			"cost_paid":  totalCost,
			"space_used": usedCapacity + spaceNeeded,
			"space_max":  maxCapacity,
		})
	}
}

func GetArmyCatalog(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		catalog, err := queries.GetTroopCatalog(r.Context())
		if err != nil {
			http.Error(w, "Failed to fetch troop catalog", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"troops_available": catalog,
		})
	}
}

type ArmyStatusResponse struct {
	TroopType    string `json:"troop_type"`
	CurrentLevel int32  `json:"current_level"`
	Quantity     int32  `json:"quantity"`
}

func GetArmyStatus(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		playerIDStr := r.Context().Value(middleware.PlayerIDKey).(string)
		parsedPlayerUUID, _ := uuid.Parse(playerIDStr)
		pgPlayerID := pgtype.UUID{Bytes: parsedPlayerUUID, Valid: true}

		rows, err := queries.GetArmyStatus(r.Context(), pgPlayerID)
		if err != nil {
			http.Error(w, "Failed to fetch army status", http.StatusInternalServerError)
			return
		}

		var army []ArmyStatusResponse
		for _, row := range rows {
			army = append(army, ArmyStatusResponse{
				TroopType:    row.TroopType,
				CurrentLevel: row.CurrentLevel,
				Quantity:     row.Quantity,
			})
		}

		if army == nil {
			army = []ArmyStatusResponse{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"army": army,
		})
	}
}
