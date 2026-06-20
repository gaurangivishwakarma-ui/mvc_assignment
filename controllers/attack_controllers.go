package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeployedTroop struct {
	TroopType string `json:"troop_type"`
	Level     int32  `json:"level"`
	Quantity  int32  `json:"quantity"`
}

type InstantBattleRequest struct {
	OpponentID     string          `json:"opponent_id"`
	DeployedTroops []DeployedTroop `json:"deployed_troops"`
}

func AttackOpponent(pool *pgxpool.Pool) http.HandlerFunc {
	queries := db.New(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		playerIDStr := r.Context().Value(middleware.PlayerIDKey).(string)
		parsedAttackerUUID, _ := uuid.Parse(playerIDStr)
		pgAttackerID := pgtype.UUID{Bytes: parsedAttackerUUID, Valid: true}

		var req InstantBattleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		parsedOpponentUUID, err := uuid.Parse(req.OpponentID)
		if err != nil {
			http.Error(w, "Invalid opponent ID format", http.StatusBadRequest)
			return
		}
		pgOpponentID := pgtype.UUID{Bytes: parsedOpponentUUID, Valid: true}

		var troopTypes []string
		var levels []int32
		var quantities []int32

		for _, t := range req.DeployedTroops {
			if t.Level <= 0 {
				t.Level = 1
			}
			troopTypes = append(troopTypes, t.TroopType)
			levels = append(levels, t.Level)
			quantities = append(quantities, t.Quantity)
		}

		attackerPower, err := queries.GetArmyCombatPower(r.Context(), db.GetArmyCombatPowerParams{
			TroopTypes: troopTypes,
			Levels:     levels,
			Quantities: quantities,
		})
		if err != nil {
			log.Printf("Database Error calculating attacker power: %v\n", err)
			http.Error(w, "Failed to compute attacker strength metrics", http.StatusInternalServerError)
			return
		}

		defensePower, err := queries.GetDefenderDefensePower(r.Context(), pgOpponentID)
		if err != nil {
			log.Printf("Database Error calculating defense power: %v\n", err)
			http.Error(w, "Failed to compute defender structure layout", http.StatusInternalServerError)
			return
		}

		var destructionPercent int32 = 100
		if defensePower > 0 {
			destructionPercent = (attackerPower * 100) / defensePower
		}
		if destructionPercent > 100 {
			destructionPercent = 100
		}

		var isAttackerWinner bool
		var xpChange int32 = 0
		var lootPercentage int32 = 0

		switch {
		case destructionPercent < 30:
			isAttackerWinner = false
			xpChange = 10
			lootPercentage = 0
		case destructionPercent < 50:
			isAttackerWinner = false
			xpChange = 5
			lootPercentage = 10
		case destructionPercent < 75:
			isAttackerWinner = true
			xpChange = 10
			lootPercentage = 20
		case destructionPercent < 100:
			isAttackerWinner = true
			xpChange = 20
			lootPercentage = 30
		default:
			isAttackerWinner = true
			xpChange = 30
			lootPercentage = 40
		}

		defender, err := queries.GetPlayerProfile(r.Context(), pgOpponentID)
		if err != nil {
			http.Error(w, "Opponent profile not found", http.StatusNotFound)
			return
		}

		stolenGold := (defender.GoldCoins * lootPercentage) / 100
		stolenElixir := (defender.Elixir * lootPercentage) / 100

		tx, err := pool.Begin(r.Context())
		if err != nil {
			http.Error(w, "Database transaction failed to initialize", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(r.Context())

		qtx := queries.WithTx(tx)

		for _, t := range req.DeployedTroops {
			if t.Level <= 0 {
				t.Level = 1
			}

			ownedQty, err := qtx.GetPlayerTroopQuantity(r.Context(), db.GetPlayerTroopQuantityParams{
				PlayerID:     pgAttackerID,
				TroopType:    t.TroopType,
				CurrentLevel: t.Level,
			})
			if err != nil {
				errMsg := fmt.Sprintf("You do not have any trained troops of type: %s at Level %d", t.TroopType, t.Level)
				http.Error(w, errMsg, http.StatusBadRequest)
				return
			}

			if ownedQty < int32(t.Quantity) {
				http.Error(w, "Insufficient army composition! You tried deploying more troops than you have trained.", http.StatusBadRequest)
				return
			}

			err = qtx.DeductPlayerTroops(r.Context(), db.DeductPlayerTroopsParams{
				PlayerID:     pgAttackerID,
				TroopType:    t.TroopType,
				CurrentLevel: t.Level,
				Quantity:     t.Quantity,
			})
			if err != nil {
				http.Error(w, "Database transaction failed during troop deduction", http.StatusInternalServerError)
				return
			}
		}

		if isAttackerWinner {
			if err := qtx.DeductLootAndXPFromDefender(r.Context(), db.DeductLootAndXPFromDefenderParams{
				GoldCoins: stolenGold, Elixir: stolenElixir, XpPoints: xpChange, ID: pgOpponentID,
			}); err != nil {
				return
			}
			if err := qtx.AddLootAndXPToAttacker(r.Context(), db.AddLootAndXPToAttackerParams{
				GoldCoins: stolenGold, Elixir: stolenElixir, XpPoints: xpChange, ID: pgAttackerID,
			}); err != nil {
				return
			}
		} else {
			if err := qtx.DeductLootAndXPFromDefender(r.Context(), db.DeductLootAndXPFromDefenderParams{
				GoldCoins: stolenGold, Elixir: stolenElixir, XpPoints: 0, ID: pgOpponentID,
			}); err != nil {
				return
			}
			if err := qtx.AddLootAndXPToAttacker(r.Context(), db.AddLootAndXPToAttackerParams{
				GoldCoins: stolenGold, Elixir: stolenElixir, XpPoints: 0, ID: pgAttackerID,
			}); err != nil {
				return
			}
			if err := qtx.DeductLootAndXPFromDefender(r.Context(), db.DeductLootAndXPFromDefenderParams{
				GoldCoins: 0, Elixir: 0, XpPoints: xpChange, ID: pgAttackerID,
			}); err != nil {
				return
			}
			if err := qtx.AddLootAndXPToAttacker(r.Context(), db.AddLootAndXPToAttackerParams{
				GoldCoins: 0, Elixir: 0, XpPoints: xpChange, ID: pgOpponentID,
			}); err != nil {
				return
			}
		}

		logData := map[string]interface{}{
			"deployed_troops": req.DeployedTroops,
			"summary": map[string]interface{}{
				"computed_attacker_power": attackerPower,
				"computed_defense_power":  defensePower,
				"destruction_percent":     destructionPercent,
			},
		}

		jsonBytes, err := json.Marshal(logData)
		if err != nil {
			http.Error(w, "Failed to compile battle logs", http.StatusInternalServerError)
			return
		}

		uuidNew, _ := uuid.NewRandom()
		pgBattleID := pgtype.UUID{Bytes: uuidNew, Valid: true}
		var damageNumeric pgtype.Numeric
		damageNumeric.Int = big.NewInt(int64(destructionPercent))
		damageNumeric.Valid = true

		if err := qtx.RecordBattle(r.Context(), db.RecordBattleParams{
			ID:               pgBattleID,
			AttackerID:       pgAttackerID,
			DefenderID:       pgOpponentID,
			IsAttackerWinner: pgtype.Bool{Bool: isAttackerWinner, Valid: true},
			GoldStolen:       stolenGold,
			ElixirStolen:     stolenElixir,
			DamagePercentage: damageNumeric,
			BattleLogs:       jsonBytes,
		}); err != nil {
			return
		}

		if err := tx.Commit(r.Context()); err != nil {
			http.Error(w, "Failed to commit battle data ledger", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":              "Simulation complete!",
			"victory":             isAttackerWinner,
			"destruction_percent": destructionPercent,
			"xp_modifier":         xpChange,
			"loot_stolen": map[string]int32{
				"gold_coins": stolenGold,
				"elixir":     stolenElixir,
			},
		})
	}
}
