-- name: GetTroopCatalog :many
SELECT * FROM troops;

-- name: GetTroopDetails :one
SELECT troop_type, elixir_cost, housing_space, level_req 
FROM troops 
WHERE troop_type = $1 AND level_req = $2;

-- name: GetTotalHousingCapacity :one
SELECT COALESCE(SUM(ac.housing_capacity), 0)::int AS total_capacity
FROM buildings_owned bo
JOIN army_camp_buildings ac ON bo.building_id = ac.building_id
WHERE bo.player_id = $1;

-- name: GetCurrentHousingUsed :one
SELECT COALESCE(SUM(t_own.quantity * t.housing_space), 0)::int AS used_capacity
FROM troops_owned t_own
JOIN troops t ON t_own.troop_type = t.troop_type AND t_own.current_level = t.level_req
WHERE t_own.player_id = $1;

-- name: PayForTroopTraining :exec
UPDATE player
SET elixir = elixir - sqlc.arg('elixir')
WHERE id = sqlc.arg('id');

-- name: AddTroopsToArmy :exec
INSERT INTO troops_owned (player_id, troop_type, current_level, quantity)
VALUES ($1, $2, $3, $4)
ON CONFLICT (player_id, troop_type) 
DO UPDATE SET quantity = troops_owned.quantity + EXCLUDED.quantity;  

-- name: GetArmyStatus :many
SELECT troop_type, current_level, quantity
FROM troops_owned
WHERE player_id = $1 AND quantity > 0;