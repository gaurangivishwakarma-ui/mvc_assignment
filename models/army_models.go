package models

type TrainTroopRequest struct {
	TroopType string `json:"troop_type"`
	Level     int32  `json:"level"`
	Quantity  int32  `json:"quantity"`
}

type ArmyStatusResponse struct {
	TroopType    string `json:"troop_type"`
	CurrentLevel int32  `json:"current_level"`
	Quantity     int32  `json:"quantity"`
}
