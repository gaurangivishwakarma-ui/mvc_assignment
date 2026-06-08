-- name: CreatePlayerAuth :one
INSERT INTO authorisation (player_id, username, password_hash)
VALUES ($1, $2, $3)
RETURNING player_id, username;

-- name: CreatePlayerProfile :one
INSERT INTO player (id)
VALUES ($1)
RETURNING id, village_level, gold_coins, elixir;

-- name: GetAuthByUsername :one
SELECT player_id, username, password_hash
FROM authorisation
WHERE username = $1 LIMIT 1;
