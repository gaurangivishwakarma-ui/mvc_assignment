INSERT INTO troops (troop_type, level_req, elixir_cost, hit_points, hit_rate, damage, speed, attack_range, housing_space, description) VALUES

-- Barbarian: cheap melee, low range, high speed
('barbarian', 1, 50,   300,  1, 75,  4, 1, 1, 'A fierce warrior who loves to fight. Cheap and fast.'),
('barbarian', 2, 100,  450,  1, 110, 4, 1, 1, 'Stronger barbarian with a sharper blade.'),

-- Archer: ranged, low HP, moderate damage
('archer',    1, 100,  190,  1, 55,  3, 3, 1, 'A sharp-eyed bowwoman who can attack from a distance.'),
('archer',    2, 200,  270,  1, 80,  3, 3, 1, 'Veteran archer with a reinforced bow.'),

-- Giant: high HP tank, targets defenses only
('giant',     1, 500,  1100, 1, 70,  2, 1, 5, 'A slow but mighty giant who charges straight at defenses.'),
('giant',     2, 1000, 1600, 1, 110, 2, 1, 5, 'An enormous giant who can soak up tremendous damage.'),

-- Goblin: fast, targets resources, double damage on resource buildings
('goblin',    1, 25,   160,  1, 55,  5, 1, 1, 'A sneaky goblin obsessed with stealing gold and elixir.'),
('goblin',    2, 50,   240,  1, 85,  5, 1, 1, 'A greedier goblin who can carry even more loot.'),

-- Wizard: area damage, high elixir cost
('wizard',    1, 1500, 670,  1, 175, 3, 3, 4, 'A powerful spellcaster who deals area damage.'),
('wizard',    2, 2000, 950,  1, 240, 3, 3, 4, 'A master wizard whose fireballs level entire rows.');

INSERT INTO buildings (id, building_type, level, name, width, breadth, cost_type, cost, hit_points, level_req) VALUES

(101, 'gold_mine', 1, 'Gold Mine Lv1',      3, 3, 'gold',   150,   400, 1),
(102, 'gold_mine', 2, 'Gold Mine Lv2',      3, 3, 'gold',   300,   500, 1),
(103, 'gold_mine', 3, 'Gold Mine Lv3',      3, 3, 'gold',   700,   620, 2),
(104, 'gold_mine', 4, 'Gold Mine Lv4',      3, 3, 'gold',  1400,   760, 3),

(111, 'elixir_collector', 1, 'Elixir Collector Lv1', 3, 3, 'elixir',  150,  400, 1),
(112, 'elixir_collector', 2, 'Elixir Collector Lv2', 3, 3, 'elixir',  300,  500, 1),
(113, 'elixir_collector', 3, 'Elixir Collector Lv3', 3, 3, 'elixir',  700,  620, 2),
(114, 'elixir_collector', 4, 'Elixir Collector Lv4', 3, 3, 'elixir', 1400,  760, 3),

(201, 'gold_storage', 1, 'Gold Storage Lv1', 3, 3, 'gold',   300,   600, 1),
(202, 'gold_storage', 2, 'Gold Storage Lv2', 3, 3, 'gold',   750,   800, 1),
(203, 'gold_storage', 3, 'Gold Storage Lv3', 3, 3, 'gold',  1500,  1000, 2),
(204, 'gold_storage', 4, 'Gold Storage Lv4', 3, 3, 'gold',  3000,  1200, 3),

(211, 'elixir_storage', 1, 'Elixir Storage Lv1', 3, 3, 'elixir',  300,   600, 1),
(212, 'elixir_storage', 2, 'Elixir Storage Lv2', 3, 3, 'elixir',  750,   800, 1),
(213, 'elixir_storage', 3, 'Elixir Storage Lv3', 3, 3, 'elixir', 1500,  1000, 2),
(214, 'elixir_storage', 4, 'Elixir Storage Lv4', 3, 3, 'elixir', 3000,  1200, 3),

(301, 'army_camp', 1, 'Army Camp Lv1', 4, 4, 'gold',   200,   500, 1),
(302, 'army_camp', 2, 'Army Camp Lv2', 4, 4, 'gold',   500,   650, 1),
(303, 'army_camp', 3, 'Army Camp Lv3', 4, 4, 'gold',  1000,   800, 2),
(304, 'army_camp', 4, 'Army Camp Lv4', 4, 4, 'gold',  2000,  1000, 3),

(401, 'cannon', 1, 'Cannon Lv1', 3, 3, 'gold',   250,   420, 1),
(402, 'cannon', 2, 'Cannon Lv2', 3, 3, 'gold',   500,   580, 1),
(403, 'cannon', 3, 'Cannon Lv3', 3, 3, 'gold',  1000,   740, 2),
(404, 'cannon', 4, 'Cannon Lv4', 3, 3, 'gold',  2000,   900, 3),

(411, 'archer_tower', 1, 'Archer Tower Lv1', 3, 3, 'gold',   400,   380, 1),
(412, 'archer_tower', 2, 'Archer Tower Lv2', 3, 3, 'gold',   800,   520, 2),
(413, 'archer_tower', 3, 'Archer Tower Lv3', 3, 3, 'gold',  1500,   680, 3),
(414, 'archer_tower', 4, 'Archer Tower Lv4', 3, 3, 'gold',  3000,   850, 4),

(421, 'mortar', 1, 'Mortar Lv1', 4, 4, 'gold',  1000,   400, 2),
(422, 'mortar', 2, 'Mortar Lv2', 4, 4, 'gold',  2000,   500, 2),
(423, 'mortar', 3, 'Mortar Lv3', 4, 4, 'gold',  4000,   620, 3),
(424, 'mortar', 4, 'Mortar Lv4', 4, 4, 'gold',  7000,   760, 4);

INSERT INTO resource_buildings (building_id, production_rate, resource_type) VALUES
(101, 100.0,  'gold'),
(102, 200.0,  'gold'),
(103, 400.0,  'gold'),
(104, 700.0,  'gold'),

(111, 100.0,  'elixir'),
(112, 200.0,  'elixir'),
(113, 400.0,  'elixir'),
(114, 700.0,  'elixir');

INSERT INTO storage_buildings (building_id, storage_capacity, resource_type) VALUES
(201,  10000, 'gold'),
(202,  50000, 'gold'),
(203, 150000, 'gold'),
(204, 450000, 'gold'),

(211,  10000, 'elixir'),
(212,  50000, 'elixir'),
(213, 150000, 'elixir'),
(214, 450000, 'elixir');

INSERT INTO army_camp_buildings (building_id, housing_capacity) VALUES
(301, 20),
(302, 30),
(303, 40),
(304, 50);

-- Cannon: slow, hits hard, short range
INSERT INTO defense_buildings (building_id, damage, attack_range, attack_speed, hit_rate) VALUES
(401,  40, 5, 1.0, 1),
(402,  70, 5, 1.0, 1),
(403, 110, 6, 0.9, 1),
(404, 160, 6, 0.8, 1);

-- Archer Tower: fast, moderate damage, longer range
INSERT INTO defense_buildings (building_id, damage, attack_range, attack_speed, hit_rate) VALUES
(411,  25, 7, 0.5, 1),
(412,  40, 7, 0.5, 1),
(413,  60, 8, 0.4, 1),
(414,  90, 8, 0.4, 1);

-- Mortar: area damage, very slow, long range, unlocks at village level 2
INSERT INTO defense_buildings (building_id, damage, attack_range, attack_speed, hit_rate) VALUES
(421, 100, 10, 5.0, 1),
(422, 150, 10, 4.5, 1),
(423, 200, 11, 4.0, 1),
(424, 280, 11, 3.5, 1);