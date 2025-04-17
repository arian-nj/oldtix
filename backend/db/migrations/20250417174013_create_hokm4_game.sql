-- migrate:up
CREATE TABLE hokm4_game (
  id BIGSERIAL NOT NULL PRIMARY KEY,

  team_one_tricks_score BIGINT Not NULL, 
  team_two_tricks_score BIGINT Not NULL,
  created_stamp TIMESTAMPTZ DEFAULT NOW(),
  end_stamp TIMESTAMPTZ
);

-- migrate:down

