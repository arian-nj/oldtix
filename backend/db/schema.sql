CREATE TABLE patch_version (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  rmode TEXT NOT NULL,
  version_number TEXT NOT NULL
);

CREATE TABLE person (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  created TIMESTAMPTZ DEFAULT NOW(),
  username TEXT NOT NULL,
  display_name TEXT,
  hashed_password TEXT NOT NULL,
  bio  text,
  coin BIGINT DEFAULT 50
);

CREATE TABLE user_statistic(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT DEFAULT 0 REFERENCES person(id),
    win BIGINT DEFAULT 0,
    lose BIGINT DEFAULT 0,
    total_tricks_won BIGINT DEFAULT 0,
    total_tricks_lost BIGINT DEFAULT 0,
    total_turns_won BIGINT DEFAULT 0,
    total_turns_lost BIGINT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE hokm4_game (
  id BIGSERIAL NOT NULL PRIMARY KEY,

  team_one_tricks_score BIGINT DEFAULT 0, 
  team_two_tricks_score BIGINT DEFAULT 0,
  created_stamp TIMESTAMPTZ DEFAULT NOW(),
  end_stamp TIMESTAMPTZ
);

-- CREATE TABLE hokm4_game_statistic(
--     id BIGSERIAL NOT NULL PRIMARY KEY,
--     match_id BIGINT NOT NULL REFERENCES hokm4_game(id),
--     person_id BIGINT REFERENCES person(id),
--     tricks_won BIGINT NOT NULL,
--     tricks_lost BIGINT NOT NULL,
--     turns_won BIGINT NOT NULL,
--     turns_lost BIGINT NOT NULL,
--     is_won BOOLEAN NOT NULL,
--     created_at TIMESTAMP DEFAULT NOW()
-- );


CREATE TABLE trick (
  trick_id BIGSERIAL NOT NULL PRIMARY KEY,
  
  game_id BIGINT NOT NULL,

  hokm BIGINT,

  hakem_index BIGINT NOT NULL,

  team_one_turn_score BIGINT DEFAULT 0,
  team_two_turn_score BIGINT DEFAULT 0,
  
  CONSTRAINT fk_game FOREIGN KEY (game_id) REFERENCES hokm4_game(id)
);

CREATE TABLE turn (
  turn_id BIGSERIAL NOT NULL PRIMARY KEY,
  trick_id BIGINT NOT NULL,
  moves TEXT NOT NULL,
  CONSTRAINT fk_trick FOREIGN KEY (trick_id) REFERENCES trick(trick_id)
);

CREATE TABLE game_player (
  player_id BIGINT NOT NULL,
  game_id BIGINT NOT NULL,
  team BIGINT NOT NULL,
  PRIMARY KEY (player_id, game_id),
  CONSTRAINT fk_player Foreign Key (player_id) REFERENCES person(id),
  CONSTRAINT fk_game Foreign Key (game_id) REFERENCES hokm4_game(id)
)

-- DROP TABLE hokm4_game;