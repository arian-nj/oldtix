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
  coin integer DEFAULT 0
);

CREATE TABLE hokm4_game (
  id BIGSERIAL NOT NULL PRIMARY KEY,

  team_one_tricks_score INT NOT NULL DEFAULT 0,
  team_two_tricks_score INT NOT NULL DEFAULT 0
);

CREATE TABLE trick (
  trick_id BIGSERIAL NOT NULL PRIMARY KEY,
  
  game_id BIGSERIAL NOT NULL,

  hokm INT,

  hakem_index INT NOT NULL,

  team_one_turn_score INT NOT NULL DEFAULT 0,
  team_two_turn_score INT NOT NULL DEFAULT 0,
  
  CONSTRAINT fk_game FOREIGN KEY (game_id) REFERENCES hokm4_game(id)
);

CREATE TABLE turn (
  turn_id BIGSERIAL NOT NULL PRIMARY KEY,
  trick_id BIGSERIAL NOT NULL,
  moves TEXT NOT NULL,
  CONSTRAINT fk_trick FOREIGN KEY (trick_id) REFERENCES trick(trick_id)
);

CREATE TABLE game_player (
  player_id BIGINT NOT NULL,
  game_id BIGINT NOT NULL,
  team SMALLINT NOT NULL,
  PRIMARY KEY (player_id, game_id),
  CONSTRAINT fk_player Foreign Key (player_id) REFERENCES person(id),
  CONSTRAINT fk_game Foreign Key (game_id) REFERENCES hokm4_game(id)
)

-- DROP TABLE hokm4_game;