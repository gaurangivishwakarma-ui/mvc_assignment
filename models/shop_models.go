package models

type PurchaseBuildingRequest struct {
	BuildingID int32 `json:"building_id"`
	XCoords    int32 `json:"x_coords"`
	YCoords    int32 `json:"y_coords"`
}

type CompleteBuildRequest struct {
	PlacementID string `json:"placement_id"`
}
