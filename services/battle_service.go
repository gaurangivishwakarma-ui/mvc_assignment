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

func GetMatch(ctx context.Context, queries *db.Queries, pgPlayerID pgtype.UUID) (map[string]interface{}, int, error) {
	attacker, err := queries.GetPlayerProfile(ctx, pgPlayerID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to load attacker profile")
	}

	matchArgs := db.GetSuitableOpponentParams{
		AttackerID:   pgPlayerID,
		VillageLevel: attacker.VillageLevel,
		XpPoints:     attacker.XpPoints,
	}

	opponent, err := queries.GetSuitableOpponent(ctx, matchArgs)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("No suitable opponents found right now. Try again later!")
	}

	rows, _ := queries.GetOpponentVillage(ctx, opponent.OpponentID)

	var village []models.OpponentBuilding
	for _, row := range rows {
		var uidStr string
		if parsed, err := uuid.FromBytes(row.PlacementID.Bytes[:]); err == nil {
			uidStr = parsed.String()
		}

		village = append(village, models.OpponentBuilding{
			PlacementID:  uidStr,
			BuildingID:   row.BuildingID,
			BuildingType: row.BuildingType,
			CurrentLevel: row.CurrentLevel,
			XCoords:      row.XCoords,
			YCoords:      row.YCoords,
		})
	}
	if village == nil {
		village = []models.OpponentBuilding{}
	}

	lootableGold := opponent.GoldCoins
	lootableElixir := opponent.Elixir

	var oppIdStr string
	if parsed, err := uuid.FromBytes(opponent.OpponentID.Bytes[:]); err == nil {
		oppIdStr = parsed.String()
	}

	return map[string]interface{}{
		"opponent_id":   oppIdStr,
		"username":      opponent.Username,
		"village_level": opponent.VillageLevel,
		"xp_points":     opponent.XpPoints,
		"loot_available": map[string]int32{
			"gold_coins": lootableGold,
			"elixir":     lootableElixir,
		},
		"village_layout": village,
	}, http.StatusOK, nil
}

func GetBattleReplay(ctx context.Context, queries *db.Queries, pgPlayerID pgtype.UUID, battleIDStr string) (models.BattleReplayResponse, int, error) {
	parsedBattleUUID, err := uuid.Parse(battleIDStr)
	if err != nil {
		return models.BattleReplayResponse{}, http.StatusBadRequest, fmt.Errorf("Invalid battle ID format")
	}
	pgBattleID := pgtype.UUID{Bytes: parsedBattleUUID, Valid: true}

	battle, err := queries.GetBattleByID(ctx, pgBattleID)
	if err != nil {
		return models.BattleReplayResponse{}, http.StatusNotFound, fmt.Errorf("Battle record not found")
	}

	if battle.AttackerID != pgPlayerID && battle.DefenderID != pgPlayerID {
		return models.BattleReplayResponse{}, http.StatusForbidden, fmt.Errorf("Unauthorized: You did not participate in this battle")
	}

	var attackerUUID, defenderUUID uuid.UUID
	copy(attackerUUID[:], battle.AttackerID.Bytes[:])
	copy(defenderUUID[:], battle.DefenderID.Bytes[:])

	var damageFloat float64
	if battle.DamagePercentage.Valid {
		_ = battle.DamagePercentage.Scan(&damageFloat)
	}

	res := models.BattleReplayResponse{
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
		res.BattleLogs = battle.BattleLogs
	}

	return res, http.StatusOK, nil
}
