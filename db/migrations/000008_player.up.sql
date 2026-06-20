CREATE TABLE "player" (
  "id" uuid PRIMARY KEY,
  "village_level" int NOT NULL CHECK (village_level >= 1) DEFAULT 1,
  "gold_coins" int NOT NULL CHECK (gold_coins >= 0) DEFAULT 1000,
  "elixir" int NOT NULL CHECK (elixir >= 0) DEFAULT 500,
  "xp_points" int NOT NULL CHECK (xp_points >= 0) DEFAULT 0,
  "attacks_won" int NOT NULL CHECK (attacks_won >= 0) DEFAULT 0,
  "attacks_lost" int NOT NULL CHECK (attacks_lost >= 0) DEFAULT 0,
  "total_looted" int NOT NULL CHECK (total_looted >= 0) DEFAULT 0,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "is_deleted" bool NOT NULL DEFAULT false
);

CREATE INDEX ON "player" ("xp_points");

ALTER TABLE "player" ADD FOREIGN KEY ("id") REFERENCES "authorisation" ("player_id") DEFERRABLE INITIALLY IMMEDIATE;
