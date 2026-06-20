CREATE TABLE "resource_buildings" (
  "building_id" int PRIMARY KEY,
  "production_rate" decimal NOT NULL CHECK (production_rate > 0),
  "resource_type" resource_type NOT NULL
);

ALTER TABLE "resource_buildings" ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") DEFERRABLE INITIALLY IMMEDIATE;
