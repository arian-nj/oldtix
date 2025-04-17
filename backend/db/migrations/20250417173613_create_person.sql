-- migrate:up
CREATE TABLE person (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  created TIMESTAMPTZ DEFAULT NOW(),
  username TEXT NOT NULL,
  display_name TEXT,
  hashed_password TEXT NOT NULL,
  bio  text,
  coin BIGINT Not NULL
);


-- migrate:down
