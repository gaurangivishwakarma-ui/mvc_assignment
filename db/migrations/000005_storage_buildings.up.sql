CREATE TABLE "storage_buildings" (
  "building_id" int PRIMARY KEY,
  "storage_capacity" int NOT NULL CHECK (storage_capacity > 0),
  "resource_type" resource_type NOT NULL
);

ALTER TABLE "storage_buildings" ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") DEFERRABLE INITIALLY IMMEDIATE;
