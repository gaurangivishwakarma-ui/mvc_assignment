-- name: GetOpponentVillage :many
SELECT 
    bo.id AS placement_id, 
    bo.building_id, 
    bo.building_type, 
    bo.current_level, 
    bo.x_coords, 
    bo.y_coords,
    b.width,
    b.breadth
FROM buildings_owned bo
JOIN buildings b ON bo.building_id = b.id
WHERE bo.player_id = $1;

-- name: GetSuitableOpponent :one
SELECT 
    p.id AS opponent_id, 
    a.username, 
    p.village_level, 
    p.gold_coins, 
    p.elixir,
    p.xp_points
FROM player p
JOIN authorisation a ON p.id = a.player_id
WHERE p.id != sqlc.arg(attacker_id)
  AND ABS(p.village_level - sqlc.arg(village_level)::int) <= 1
  AND ABS(p.xp_points - sqlc.arg(xp_points)::int) <= 500
ORDER BY RANDOM()
LIMIT 1;

-- name: GetArmyCombatPower :one
WITH deployed AS (
    SELECT unnest(sqlc.arg(troop_types)::varchar[]) AS troop_type,
           unnest(sqlc.arg(levels)::int[]) AS troop_level,
           unnest(sqlc.arg(quantities)::int[]) AS qty
)
SELECT 
    COALESCE(SUM((t.damage + t.hit_points) * d.qty), 0)::int AS attacker_power
FROM deployed d
JOIN troops t ON t.troop_type = d.troop_type AND t.level_req = d.troop_level;

-- name: GetDefenderDefensePower :one
SELECT 
    COALESCE(SUM(b.hit_points + COALESCE(db.damage, 0)), 0)::int AS defense_power
FROM buildings_owned bo
JOIN buildings b ON bo.building_id = b.id
LEFT JOIN defense_buildings db ON b.id = db.building_id
WHERE bo.player_id = $1;

-- name: DeductLootAndXPFromDefender :exec
UPDATE player 
SET gold_coins = GREATEST(0, gold_coins - $1), 
    elixir = GREATEST(0, elixir - $2),
    xp_points = GREATEST(0, xp_points - $3)
WHERE id = $4;

-- name: AddLootAndXPToAttacker :exec
UPDATE player 
SET gold_coins = gold_coins + $1, 
    elixir = elixir + $2,
    xp_points = xp_points + $3
WHERE id = $4;

-- name: RecordBattle :exec
INSERT INTO battles (
    id, 
    attacker_id, 
    defender_id, 
    is_attacker_winner, 
    gold_stolen, 
    elixir_stolen, 
    damage_percentage, 
    battle_logs 
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8); 

-- name: GetPlayerTroopQuantity :one
SELECT quantity 
FROM troops_owned 
WHERE player_id = sqlc.arg('player_id') 
  AND troop_type = sqlc.arg('troop_type')
  AND current_level = sqlc.arg('current_level');


-- name: DeductPlayerTroops :exec
UPDATE troops_owned 
SET quantity = quantity - sqlc.arg('quantity')
WHERE player_id = sqlc.arg('player_id') 
  AND troop_type = sqlc.arg('troop_type')
  AND current_level = sqlc.arg('current_level');

-- name: DeletePlayerTroops :exec
DELETE FROM troops_owned
WHERE player_id = sqlc.arg('player_id') 
  AND troop_type = sqlc.arg('troop_type')
  AND current_level = sqlc.arg('current_level');


-- name: GetBattleByID :one
SELECT 
    id, 
    attacker_id, 
    defender_id, 
    is_attacker_winner, 
    gold_stolen, 
    elixir_stolen, 
    damage_percentage, 
    battle_logs, 
    battle_time
FROM battles
WHERE id = $1 LIMIT 1;
