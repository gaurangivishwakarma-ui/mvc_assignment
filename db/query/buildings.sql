-- name: GetOwnedBuildingPositionInfo :one
SELECT bo.player_id, b.width, b.breadth
FROM buildings_owned bo
JOIN buildings b ON bo.building_id = b.id
WHERE bo.id = $1;

-- name: UpdateBuildingPosition :exec
UPDATE buildings_owned
SET x_coords = $2, y_coords = $3
WHERE id = $1 AND player_id = $4;

-- name: GetShopCatalog :many
SELECT * FROM buildings 
WHERE level = 1;