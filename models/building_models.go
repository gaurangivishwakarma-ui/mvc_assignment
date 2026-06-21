package models

type MoveBuildingRequest struct {
	OwnedBuildingID string `json:"owned_building_id"`
	NewX            int32  `json:"new_x"`
	NewY            int32  `json:"new_y"`
}
