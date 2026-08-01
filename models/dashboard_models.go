package models

type BuildingInventory struct {
	PlacementID  string `json:"placement_id,omitempty"`
	BuildingID   int32  `json:"building_id,omitempty"`
	BuildingType string `json:"building_type,omitempty"`
	CurrentLevel int32  `json:"current_level,omitempty"`
	XCoords      int32  `json:"x_coords"`
	YCoords      int32  `json:"y_coords"`
	IsBuilt      bool   `json:"is_built"`
}

type DashboardResponse struct {
	Username          string              `json:"username"`
	VillageLevel      int32               `json:"village_level"`
	XPPoints          int32               `json:"xp_points"`
	Balances          map[string]int32    `json:"balances"`
	StorageCapacities map[string]int32    `json:"storage_capacities"`
	Buildings         []BuildingInventory `json:"buildings_placed"`
}
