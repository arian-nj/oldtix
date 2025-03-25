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
  coin bigserial NOT NULL
);

CREATE TABLE user_statistic(
    user_id BIGSERIAL PRIMARY KEY REFERENCES person(id),
    wins bigserial NOT NULL,
    losses bigserial NOT NULL,
    total_tricks_won bigserial NOT NULL,
    total_tricks_lost bigserial NOT NULL,
    total_turns_won bigserial NOT NULL,
    total_turns_lost bigserial NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION create_user_statistic()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_statistic(
      user_id, wins, losses, total_tricks_won, total_tricks_lost, total_turns_won, total_turns_lost
    ) VALUES (
      NEW.id, 0, 0, 0, 0, 0, 0
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_user_statistic
AFTER INSERT ON person
FOR EACH ROW
EXECUTE PROCEDURE create_user_statistic();

CREATE TABLE hokm4_game (
  id BIGSERIAL NOT NULL PRIMARY KEY,

  team_one_tricks_score bigserial NOT NULL,
  team_two_tricks_score bigserial NOT NULL,
  created_stamp TIMESTAMPTZ DEFAULT NOW(),
  end_stamp TIMESTAMPTZ
);

CREATE TABLE hokm4_game_statistic(
    match_id BIGSERIAL PRIMARY KEY REFERENCES hokm4_game(id),
    person_id BIGSERIAL REFERENCES person(id),
    tricks_won bigserial NOT NULL,
    tricks_lost bigserial NOT NULL,
    turns_won bigserial NOT NULL,
    turns_lost bigserial NOT NULL,
    is_won BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE trick (
  trick_id BIGSERIAL NOT NULL PRIMARY KEY,
  
  game_id BIGSERIAL NOT NULL,

  hokm bigserial,

  hakem_index bigserial NOT NULL,

  team_one_turn_score bigserial NOT NULL,
  team_two_turn_score bigserial NOT NULL,
  
  CONSTRAINT fk_game FOREIGN KEY (game_id) REFERENCES hokm4_game(id)
);

CREATE TABLE turn (
  turn_id BIGSERIAL NOT NULL PRIMARY KEY,
  trick_id BIGSERIAL NOT NULL,
  moves TEXT NOT NULL,
  CONSTRAINT fk_trick FOREIGN KEY (trick_id) REFERENCES trick(trick_id)
);

CREATE TABLE game_player (
  player_id BIGSERIAL NOT NULL,
  game_id BIGSERIAL NOT NULL,
  team bigserial NOT NULL,
  PRIMARY KEY (player_id, game_id),
  CONSTRAINT fk_player Foreign Key (player_id) REFERENCES person(id),
  CONSTRAINT fk_game Foreign Key (game_id) REFERENCES hokm4_game(id)
)

-- DROP TABLE hokm4_game;