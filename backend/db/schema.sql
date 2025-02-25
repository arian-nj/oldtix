CREATE TABLE users (
  id SERIAL NOT NULL PRIMARY KEY,
  created TIMESTAMPTZ DEFAULT NOW(),
  username TEXT NOT NULL,
  display_name TEXT,
  hashed_password TEXT NOT NULL,
  bio  text,
  coin integer DEFAULT 0
);

CREATE TABLE patch_version (
  id SERIAL NOT NULL PRIMARY KEY,
  rmode TEXT NOT NULL,
  version_number TEXT NOT NULL
);

-- CREATE TABLE hokm4_game (
--   id SERIAL NOT NULL PRIMARY KEY,
--   username TEXT NOT NULL,
--   display_name TEXT,
--   hashed_password TEXT NOT NULL,
--   bio  text,
--   coin integer DEFAULT 0,
--   created TIMESTAMPTZ DEFAULT NOW()
-- );

-- DROP TABLE users;