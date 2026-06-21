package models

type DeployedTroop struct {
	TroopType string `json:"troop_type"`
	Level     int32  `json:"level"`
	Quantity  int32  `json:"quantity"`
}

type InstantBattleRequest struct {
	OpponentID     string          `json:"opponent_id"`
	DeployedTroops []DeployedTroop `json:"deployed_troops"`
}
