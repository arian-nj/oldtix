-- migrate:up
CREATE TABLE game_player (
  player_id BIGINT NOT NULL,
  game_id BIGINT NOT NULL,
  team BIGINT NOT NULL,
  PRIMARY KEY (player_id, game_id),
  CONSTRAINT fk_player Foreign Key (player_id) REFERENCES person(id),
  CONSTRAINT fk_game Foreign Key (game_id) REFERENCES hokm4_game(id)
)

-- migrate:down

