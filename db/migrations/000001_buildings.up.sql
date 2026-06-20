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

CREATE UNIQUE INDEX ON "buildings" ("building_type", "level");
