CREATE TYPE "resource_type" AS ENUM (
  'gold',
  'elixir'
);

CREATE TABLE "buildings" (
  "id" int PRIMARY KEY,
  "building_type" varchar NOT NULL,
  "level" int NOT NULL CHECK (level >= 1),
  "name" varchar NOT NULL,
  "width" int NOT NULL CHECK (width > 0),
  "breadth" int NOT NULL CHECK (breadth >0),
  "cost_type" resource_type NOT NULL,
  "cost" int NOT NULL CHECK (cost >= 0),
  "hit_points" int NOT NULL CHECK (hit_points > 0),
  "level_req" int NOT NULL CHECK (level_req >= 1)
);

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

CREATE TABLE "resource_buildings" (
  "building_id" int PRIMARY KEY,
  "production_rate" decimal NOT NULL CHECK (production_rate > 0),
  "resource_type" resource_type NOT NULL
);

CREATE TABLE "storage_buildings" (
  "building_id" int PRIMARY KEY,
  "storage_capacity" int NOT NULL CHECK (storage_capacity > 0),
  "resource_type" resource_type NOT NULL
);

CREATE TABLE "army_camp_buildings" (
  "building_id" int PRIMARY KEY,
  "housing_capacity" int NOT NULL CHECK (housing_capacity > 0)
);

CREATE TABLE "defense_buildings" (
  "building_id" int PRIMARY KEY,
  "damage" int NOT NULL CHECK (damage > 0),
  "attack_range" int NOT NULL CHECK (attack_range > 0),
  "attack_speed" decimal NOT NULL CHECK (attack_speed > 0),
  "hit_rate" int NOT NULL CHECK (hit_rate > 0)
);

CREATE TABLE "player" (
  "id" uuid PRIMARY KEY,
  "village_level" int NOT NULL CHECK (village_level >= 1) DEFAULT 1,
  "gold_coins" int NOT NULL CHECK (gold_coins >= 0) DEFAULT 0,
  "elixir" int NOT NULL CHECK (elixir >= 0) DEFAULT 0,
  "xp_points" int NOT NULL CHECK (xp_points >= 0) DEFAULT 0,
  "attacks_won" int NOT NULL CHECK (attacks_won >= 0) DEFAULT 0,
  "attacks_lost" int NOT NULL CHECK (attacks_lost >= 0) DEFAULT 0,
  "total_looted" int NOT NULL CHECK (total_looted >= 0) DEFAULT 0,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "is_deleted" bool NOT NULL DEFAULT false
);

CREATE TABLE "authorisation" (
  "player_id" uuid PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL CHECK (length(username) >= 3),
  "password_hash" varchar NOT NULL
);

CREATE TABLE "buildings_owned" (
  "id" uuid PRIMARY KEY,
  "player_id" uuid NOT NULL,
  "building_type" varchar NOT NULL,
  "building_id" int NOT NULL,
  "current_level" int NOT NULL CHECK (current_level >= 1),
  "x_coords" int NOT NULL CHECK (x_coords >= 0),
  "y_coords" int NOT NULL CHECK (y_coords >= 0),
  "last_collected_at" timestamp,
  "time_purchased" timestamp,
  "is_built" bool
);

CREATE TABLE "troops_owned" (
  "player_id" uuid NOT NULL,
  "troop_type" varchar NOT NULL,
  "current_level" int NOT NULL CHECK (current_level >= 1),
  "quantity" int NOT NULL CHECK (quantity > 0),
  PRIMARY KEY ("player_id", "troop_type")
);

CREATE TABLE "battles" (
  "id" uuid PRIMARY KEY,
  "attacker_id" uuid NOT NULL,
  "defender_id" uuid NOT NULL,
  "is_attacker_winner" bool,
  "gold_stolen" int NOT NULL CHECK (gold_stolen >= 0) DEFAULT 0,
  "elixir_stolen" int NOT NULL CHECK (elixir_stolen >= 0) DEFAULT 0,
  "damage_percentage" decimal NOT NULL CHECK (damage_percentage BETWEEN 0 AND 100),
  "battle_logs" json,
  "battle_time" timestamp NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "buildings" ("building_type", "level");

CREATE INDEX ON "player" ("xp_points");

CREATE INDEX ON "buildings_owned" ("player_id");

CREATE UNIQUE INDEX ON "buildings_owned" ("player_id", "x_coords", "y_coords");

CREATE INDEX ON "troops_owned" ("troop_type", "current_level");

CREATE INDEX ON "battles" ("attacker_id", "battle_time");

CREATE INDEX ON "battles" ("defender_id", "battle_time");

CREATE INDEX ON "battles" ("attacker_id", "is_attacker_winner");

CREATE INDEX ON "battles" ("defender_id", "is_attacker_winner");

ALTER TABLE "defense_buildings" ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "resource_buildings" ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "storage_buildings" ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "army_camp_buildings" ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "player" ADD FOREIGN KEY ("id") REFERENCES "authorisation" ("player_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "buildings_owned" ADD FOREIGN KEY ("player_id") REFERENCES "player" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "troops_owned" ADD FOREIGN KEY ("player_id") REFERENCES "player" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "battles" ADD FOREIGN KEY ("attacker_id") REFERENCES "player" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "battles" ADD FOREIGN KEY ("defender_id") REFERENCES "player" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "buildings_owned" ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "troops_owned" ADD FOREIGN KEY ("troop_type", "current_level") REFERENCES "troops" ("troop_type", "level_req") DEFERRABLE INITIALLY IMMEDIATE;