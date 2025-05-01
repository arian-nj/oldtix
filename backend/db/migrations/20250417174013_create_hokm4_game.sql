-- migrate:up
CREATE TABLE hokm4_game (
  id BIGSERIAL NOT NULL PRIMARY KEY,

  team_one_tricks_score BIGINT Not NULL DEFAULT 0, 
  team_two_tricks_score BIGINT Not NULL DEFAULT 0,
  bet_amount BIGINT NOT NULL DEFAULT 0,
  
  created_stamp TIMESTAMPTZ DEFAULT NOW(),
  end_stamp TIMESTAMPTZ
);

-- migrate:down

