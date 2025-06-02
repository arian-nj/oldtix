-- migrate:up
CREATE TABLE bot_user (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  tg_id TEXT NOT NULL UNIQUE
);
-- migrate:down
