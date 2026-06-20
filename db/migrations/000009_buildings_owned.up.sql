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

CREATE INDEX ON "buildings_owned" ("player_id");

CREATE UNIQUE INDEX ON "buildings_owned" ("player_id", "x_coords", "y_coords");

ALTER TABLE "buildings_owned" ADD FOREIGN KEY ("player_id") REFERENCES "player" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "buildings_owned" ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") DEFERRABLE INITIALLY IMMEDIATE;
