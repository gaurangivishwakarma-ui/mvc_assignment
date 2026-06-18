-- name: GetDashboardData :many
SELECT 
    a.username, 
    p.gold_coins, 
    p.elixir, 
    p.village_level, 
    p.xp_points,
    bo.id AS placement_id,
    bo.building_id,
    bo.building_type,
    bo.current_level,
    bo.x_coords,
    bo.y_coords,
    bo.is_built
FROM player p
JOIN authorisation a ON p.id = a.player_id
LEFT JOIN buildings_owned bo ON p.id = bo.player_id
WHERE p.id = $1;