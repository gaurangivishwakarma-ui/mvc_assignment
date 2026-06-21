-- name: GetPlayerProfile :one
SELECT * FROM "player"
WHERE id = $1 LIMIT 1;

-- name: GetPlayerBuildings :many
SELECT * FROM "buildings_owned"
WHERE player_id = $1;

-- name: GetPlayerBuildingsWithDimensions :many
SELECT 
    bo.id, 
    bo.player_id, 
    bo.building_type, 
    bo.building_id, 
    bo.current_level, 
    bo.x_coords, 
    bo.y_coords, 
    bo.last_collected_at, 
    bo.time_purchased, 
    bo.is_built,
    b.width,
    b.breadth
FROM buildings_owned bo
JOIN buildings b ON bo.building_id = b.id
WHERE bo.player_id = $1;

-- name: GetBuildingByID :one
SELECT * FROM buildings 
WHERE id = $1 LIMIT 1;

-- name: DeductResources :one
UPDATE player
SET gold_coins = gold_coins - $1,
    elixir = elixir - $2
WHERE id = $3 AND gold_coins >= $1 AND elixir >= $2
RETURNING id, village_level, gold_coins, elixir, xp_points;

-- name: PlaceBuilding :one
INSERT INTO buildings_owned (
    id,
    player_id, 
    building_type, 
    building_id, 
    current_level, 
    x_coords, 
    y_coords, 
    time_purchased, 
    is_built
) VALUES (
    $1, $2, $3, $4, 1, $5, $6, NOW(), false
) RETURNING *;

-- name: MarkBuildingAsBuilt :execrows
UPDATE buildings_owned
SET is_built = true
WHERE id = $1
  AND player_id = $2
  AND is_built = false
  AND NOW() >= time_purchased + INTERVAL '5 seconds';