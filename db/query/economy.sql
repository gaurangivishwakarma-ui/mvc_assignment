-- name: GetResourceBuildingsForCollection :many
SELECT 
    bo.id AS owned_id,
    bo.last_collected_at,
    b.name,
    rb.resource_type,
    CAST(rb.production_rate AS float8) AS production_rate_float  --sqlc converts float8 to float64 in go
FROM buildings_owned bo
JOIN buildings b ON bo.building_id = b.id
JOIN resource_buildings rb ON b.id = rb.building_id
WHERE bo.player_id = $1;

-- name: AddPlayerResources :one
UPDATE player
SET gold_coins = gold_coins + $1,
    elixir = elixir + $2
WHERE id = $3
RETURNING id, village_level, gold_coins, elixir, xp_points;

-- name: ResetCollectionTime :exec
UPDATE buildings_owned
SET last_collected_at = NOW()
WHERE player_id = $1 AND building_id IN (
    SELECT building_id FROM resource_buildings
);

-- name: GetTotalStorageCapacity :many
SELECT
    sb.resource_type,
    CAST(SUM(sb.storage_capacity) AS int) AS total_capacity
FROM buildings_owned bo
JOIN storage_buildings sb ON bo.building_id = sb.building_id
WHERE bo.player_id = $1
  AND bo.is_built = true
GROUP BY sb.resource_type;