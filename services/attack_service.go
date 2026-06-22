package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	db "github.com/gaurangi/mvc_assignment/db/sqlc"
	"github.com/gaurangi/mvc_assignment/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AttackOpponent(ctx context.Context, pool *pgxpool.Pool, pgAttackerID pgtype.UUID, req models.InstantBattleRequest) (map[string]interface{}, int, error) {
	queries := db.New(pool)

	parsedOpponentUUID, err := uuid.Parse(req.OpponentID)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid opponent ID format")
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

	attackerPower, err := queries.GetArmyCombatPower(ctx, db.GetArmyCombatPowerParams{
		TroopTypes: troopTypes,
		Levels:     levels,
		Quantities: quantities,
	})
	if err != nil {
		log.Printf("Database Error calculating attacker power: %v\n", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to compute attacker strength metrics")
	}

	defensePower, err := queries.GetDefenderDefensePower(ctx, pgOpponentID)
	if err != nil {
		log.Printf("Database Error calculating defense power: %v\n", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to compute defender structure layout")
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

	defender, err := queries.GetPlayerProfile(ctx, pgOpponentID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Opponent profile not found")
	}

	stolenGold := (defender.GoldCoins * lootPercentage) / 100
	stolenElixir := (defender.Elixir * lootPercentage) / 100

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Database transaction failed to initialize")
	}
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)

	for _, t := range req.DeployedTroops {
		if t.Level <= 0 {
			t.Level = 1
		}

		ownedQty, err := qtx.GetPlayerTroopQuantity(ctx, db.GetPlayerTroopQuantityParams{
			PlayerID:     pgAttackerID,
			TroopType:    t.TroopType,
			CurrentLevel: t.Level,
		})
		if err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("You do not have any trained troops of type: %s at Level %d", t.TroopType, t.Level)
		}

		if ownedQty < int32(t.Quantity) {
			return nil, http.StatusBadRequest, fmt.Errorf("Insufficient army composition! You tried deploying more troops than you have trained.")
		}

		if ownedQty == int32(t.Quantity) {
			err = qtx.DeletePlayerTroops(ctx, db.DeletePlayerTroopsParams{
				PlayerID:     pgAttackerID,
				TroopType:    t.TroopType,
				CurrentLevel: t.Level,
			})
		} else {
			err = qtx.DeductPlayerTroops(ctx, db.DeductPlayerTroopsParams{
				PlayerID:     pgAttackerID,
				TroopType:    t.TroopType,
				CurrentLevel: t.Level,
				Quantity:     t.Quantity,
			})
		}
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Database transaction failed during troop deduction")
		}
	}

	if isAttackerWinner {
		if err := qtx.DeductLootAndXPFromDefender(ctx, db.DeductLootAndXPFromDefenderParams{
			GoldCoins: stolenGold, Elixir: stolenElixir, XpPoints: xpChange, ID: pgOpponentID,
		}); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Failed to deduct resources from defender")
		}
		if err := qtx.AddLootAndXPToAttacker(ctx, db.AddLootAndXPToAttackerParams{
			GoldCoins: stolenGold, Elixir: stolenElixir, XpPoints: xpChange, ID: pgAttackerID,
		}); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Failed to add resources to attacker")
		}
	} else {
		if err := qtx.DeductLootAndXPFromDefender(ctx, db.DeductLootAndXPFromDefenderParams{
			GoldCoins: stolenGold, Elixir: stolenElixir, XpPoints: 0, ID: pgOpponentID,
		}); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Failed to deduct resources from defender")
		}
		if err := qtx.AddLootAndXPToAttacker(ctx, db.AddLootAndXPToAttackerParams{
			GoldCoins: stolenGold, Elixir: stolenElixir, XpPoints: 0, ID: pgAttackerID,
		}); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Failed to add resources to attacker")
		}
		if err := qtx.DeductLootAndXPFromDefender(ctx, db.DeductLootAndXPFromDefenderParams{
			GoldCoins: 0, Elixir: 0, XpPoints: xpChange, ID: pgAttackerID,
		}); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Failed to deduct XP from attacker")
		}
		if err := qtx.AddLootAndXPToAttacker(ctx, db.AddLootAndXPToAttackerParams{
			GoldCoins: 0, Elixir: 0, XpPoints: xpChange, ID: pgOpponentID,
		}); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("Failed to add XP to defender")
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
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to compile battle logs")
	}

	uuidNew, _ := uuid.NewRandom()
	pgBattleID := pgtype.UUID{Bytes: uuidNew, Valid: true}
	var damageNumeric pgtype.Numeric
	damageNumeric.Int = big.NewInt(int64(destructionPercent))
	damageNumeric.Valid = true

	if err := qtx.RecordBattle(ctx, db.RecordBattleParams{
		ID:               pgBattleID,
		AttackerID:       pgAttackerID,
		DefenderID:       pgOpponentID,
		IsAttackerWinner: pgtype.Bool{Bool: isAttackerWinner, Valid: true},
		GoldStolen:       stolenGold,
		ElixirStolen:     stolenElixir,
		DamagePercentage: damageNumeric,
		BattleLogs:       jsonBytes,
	}); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to record battle")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to commit battle data ledger")
	}

	return map[string]interface{}{
		"status":              "Simulation complete!",
		"victory":             isAttackerWinner,
		"destruction_percent": destructionPercent,
		"xp_modifier":         xpChange,
		"loot_stolen": map[string]int32{
			"gold_coins": stolenGold,
			"elixir":     stolenElixir,
		},
	}, http.StatusOK, nil
}
