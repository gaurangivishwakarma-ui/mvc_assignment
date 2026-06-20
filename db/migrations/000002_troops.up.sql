CREATE TABLE "troops" (
  "troop_type" varchar NOT NULL,
  "level_req" int NOT NULL CHECK (level_req >= 1),
  "elixir_cost" int NOT NULL CHECK (elixir_cost >= 0),
  "hit_points" int NOT NULL CHECK (hit_points > 0),
  "hit_rate" int NOT NULL CHECK (hit_rate > 0),
  "damage" int NOT NULL CHECK (damage >= 0),
  "speed" int NOT NULL CHECK (speed > 0),
  "attack_range" int NOT NULL CHECK (attack_range > 0),
  "housing_space" int NOT NULL CHECK (housing_space > 0),
  "description" text,
  PRIMARY KEY ("troop_type", "level_req")
);
