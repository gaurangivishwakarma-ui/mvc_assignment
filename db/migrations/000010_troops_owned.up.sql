CREATE TABLE "troops_owned" (
  "player_id" uuid NOT NULL,
  "troop_type" varchar NOT NULL,
  "current_level" int NOT NULL CHECK (current_level >= 1),
  "quantity" int NOT NULL CHECK (quantity > 0),
  PRIMARY KEY ("player_id", "troop_type")
);

CREATE INDEX ON "troops_owned" ("troop_type", "current_level");

ALTER TABLE "troops_owned" ADD FOREIGN KEY ("player_id") REFERENCES "player" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "troops_owned" ADD FOREIGN KEY ("troop_type", "current_level") REFERENCES "troops" ("troop_type", "level_req") DEFERRABLE INITIALLY IMMEDIATE;
