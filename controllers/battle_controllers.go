package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type OpponentBuilding struct {
	PlacementID  string `json:"placement_id"`
	BuildingID   int32  `json:"building_id"`
	BuildingType string `json:"building_type"`
	CurrentLevel int32  `json:"current_level"`
	XCoords      int32  `json:"x_coords"`
	YCoords      int32  `json:"y_coords"`
}

func GetMatch(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		playerIDStr := r.Context().Value(middleware.PlayerIDKey).(string)
		parsedPlayerUUID, _ := uuid.Parse(playerIDStr)
		pgPlayerID := pgtype.UUID{Bytes: parsedPlayerUUID, Valid: true}

		attacker, err := queries.GetPlayerProfile(r.Context(), pgPlayerID)
		if err != nil {
			http.Error(w, "Failed to load attacker profile", http.StatusInternalServerError)
			return
		}

		matchArgs := db.GetSuitableOpponentParams{
			AttackerID:   pgPlayerID,
			VillageLevel: attacker.VillageLevel,
			XpPoints:     attacker.XpPoints,
		}

		opponent, err := queries.GetSuitableOpponent(r.Context(), matchArgs)
		if err != nil {
			http.Error(w, "No suitable opponents found right now. Try again later!", http.StatusNotFound)
			return
		}

		rows, _ := queries.GetOpponentVillage(r.Context(), opponent.OpponentID)

		var village []OpponentBuilding
		for _, row := range rows {
			var uidStr string
			if parsed, err := uuid.FromBytes(row.PlacementID.Bytes[:]); err == nil {
				uidStr = parsed.String()
			}

			village = append(village, OpponentBuilding{
				PlacementID:  uidStr,
				BuildingID:   row.BuildingID,
				BuildingType: row.BuildingType,
				CurrentLevel: row.CurrentLevel,
				XCoords:      row.XCoords,
				YCoords:      row.YCoords,
			})
		}
		if village == nil {
			village = []OpponentBuilding{}
		}

		lootableGold := opponent.GoldCoins
		lootableElixir := opponent.Elixir

		var oppIdStr string
		if parsed, err := uuid.FromBytes(opponent.OpponentID.Bytes[:]); err == nil {
			oppIdStr = parsed.String()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"opponent_id":   oppIdStr,
			"username":      opponent.Username,
			"village_level": opponent.VillageLevel,
			"xp_points":     opponent.XpPoints,
			"loot_available": map[string]int32{
				"gold_coins": lootableGold,
				"elixir":     lootableElixir,
			},
			"village_layout": village,
		})
	}
}

type BattleReplayResponse struct {
	ID               string          `json:"id"`
	AttackerID       string          `json:"attacker_id"`
	DefenderID       string          `json:"defender_id"`
	IsAttackerWinner bool            `json:"is_attacker_winner"`
	GoldStolen       int32           `json:"gold_stolen"`
	ElixirStolen     int32           `json:"elixir_stolen"`
	DamagePercentage float64         `json:"damage_percentage"`
	BattleLogs       json.RawMessage `json:"battle_logs"`
	BattleTime       time.Time       `json:"battle_time"`
}

func GetBattleReplay(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		battleIDStr := r.URL.Query().Get("id")
		if battleIDStr == "" {
			http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
			return
		}

		parsedBattleUUID, err := uuid.Parse(battleIDStr)
		if err != nil {
			http.Error(w, "Invalid battle ID format", http.StatusBadRequest)
			return
		}
		pgBattleID := pgtype.UUID{Bytes: parsedBattleUUID, Valid: true}

		battle, err := queries.GetBattleByID(r.Context(), pgBattleID)
		if err != nil {
			http.Error(w, "Battle record not found", http.StatusNotFound)
			return
		}

		playerIDStr := r.Context().Value(middleware.PlayerIDKey).(string)
		parsedPlayerUUID, _ := uuid.Parse(playerIDStr)
		pgPlayerID := pgtype.UUID{Bytes: parsedPlayerUUID, Valid: true}

		if battle.AttackerID != pgPlayerID && battle.DefenderID != pgPlayerID {
			http.Error(w, "Unauthorized: You did not participate in this battle", http.StatusForbidden)
			return
		}

		var attackerUUID, defenderUUID uuid.UUID
		copy(attackerUUID[:], battle.AttackerID.Bytes[:])
		copy(defenderUUID[:], battle.DefenderID.Bytes[:])

		var damageFloat float64
		if battle.DamagePercentage.Valid {
			_ = battle.DamagePercentage.Scan(&damageFloat)
		}

		res := BattleReplayResponse{
			ID:               battleIDStr,
			AttackerID:       attackerUUID.String(),
			DefenderID:       defenderUUID.String(),
			IsAttackerWinner: battle.IsAttackerWinner.Bool,
			GoldStolen:       battle.GoldStolen,
			ElixirStolen:     battle.ElixirStolen,
			DamagePercentage: damageFloat,
			BattleTime:       battle.BattleTime.Time,
		}

		if len(battle.BattleLogs) > 0 {
			res.BattleLogs = json.RawMessage(battle.BattleLogs)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
