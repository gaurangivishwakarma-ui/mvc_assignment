package models

import (
	"encoding/json"
	"time"
)

type OpponentBuilding struct {
	PlacementID  string `json:"placement_id"`
	BuildingID   int32  `json:"building_id"`
	BuildingType string `json:"building_type"`
	CurrentLevel int32  `json:"current_level"`
	XCoords      int32  `json:"x_coords"`
	YCoords      int32  `json:"y_coords"`
	Width        int32  `json:"width"`
	Breadth      int32  `json:"breadth"`
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
