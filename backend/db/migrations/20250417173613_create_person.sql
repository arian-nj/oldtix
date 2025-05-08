-- migrate:up
CREATE TABLE person (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  display_name TEXT,
  bio  text,
  coin BIGINT Not NULL,
  created TIMESTAMPTZ DEFAULT NOW()
);


-- migrate:down
