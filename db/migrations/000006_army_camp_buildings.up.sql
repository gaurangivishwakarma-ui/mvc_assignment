CREATE TABLE "army_camp_buildings" (
  "building_id" int PRIMARY KEY,
  "housing_capacity" int NOT NULL CHECK (housing_capacity > 0)
);

ALTER TABLE "army_camp_buildings" ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") DEFERRABLE INITIALLY IMMEDIATE;
