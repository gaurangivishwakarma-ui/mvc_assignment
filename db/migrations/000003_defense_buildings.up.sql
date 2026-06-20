CREATE TABLE "defense_buildings" (
  "building_id" int PRIMARY KEY,
  "damage" int NOT NULL CHECK (damage > 0),
  "attack_range" int NOT NULL CHECK (attack_range > 0),
  "attack_speed" decimal NOT NULL CHECK (attack_speed > 0),
  "hit_rate" int NOT NULL CHECK (hit_rate > 0)
);

ALTER TABLE "defense_buildings" ADD FOREIGN KEY ("building_id") REFERENCES "buildings" ("id") DEFERRABLE INITIALLY IMMEDIATE;
