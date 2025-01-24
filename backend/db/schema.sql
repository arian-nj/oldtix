-- Active: 1737395128745@@127.0.0.1@5432@game
CREATE TABLE users (
  id SERIAL NOT NULL PRIMARY KEY,
  created TIMESTAMPTZ DEFAULT NOW(),
  username TEXT NOT NULL,
  display_name TEXT,
  hashed_password TEXT NOT NULL,
  bio  text,
  coin integer DEFAULT 0
);

-- DROP TABLE users;