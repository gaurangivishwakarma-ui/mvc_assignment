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

CREATE INDEX ON "battles" ("attacker_id", "battle_time");

CREATE INDEX ON "battles" ("defender_id", "battle_time");

CREATE INDEX ON "battles" ("attacker_id", "is_attacker_winner");

CREATE INDEX ON "battles" ("defender_id", "is_attacker_winner");

ALTER TABLE "battles" ADD FOREIGN KEY ("attacker_id") REFERENCES "player" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "battles" ADD FOREIGN KEY ("defender_id") REFERENCES "player" ("id") DEFERRABLE INITIALLY IMMEDIATE;
