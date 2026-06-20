CREATE TABLE "authorisation" (
  "player_id" uuid PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL CHECK (length(username) >= 3),
  "password_hash" varchar NOT NULL
);
