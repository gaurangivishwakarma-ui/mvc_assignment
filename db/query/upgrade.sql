-- name: GetOwnedBuildingDetails :one
SELECT 
    bo.id AS placement_id, 
    bo.current_level, 
    bo.is_built,
    b.building_type 
FROM buildings_owned bo
JOIN buildings b ON bo.building_id = b.id
WHERE bo.id = $1 AND bo.player_id = $2;

-- name: GetNextLevelBuilding :one
SELECT 
    id AS next_building_id, 
    cost AS build_cost, 
    cost_type, 
    level_req,
    name
FROM buildings
WHERE building_type = $1 AND level = $2;

-- name: PayForUpgrade :exec
UPDATE player
SET gold_coins = gold_coins - $1,
    elixir = elixir - $2
WHERE id = $3;

-- name: CommitBuildingUpgrade :exec
UPDATE buildings_owned
SET building_id = $1, current_level = $2, is_built = false, time_purchased = NOW()
WHERE id = $3;

-- name: UpgradeVillageLevel :exec
UPDATE player
SET village_level = village_level + 1,
    gold_coins = gold_coins - sqlc.arg('gold_cost')
WHERE id = sqlc.arg('player_id');