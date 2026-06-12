-- name: GetPlayerProfile :one
SELECT * FROM "player"
WHERE id = $1 LIMIT 1;

-- name: GetPlayerBuildings :many
SELECT * FROM "buildings_owned"
WHERE player_id = $1;