-- migrate:up
CREATE TABLE trick (
  trick_id BIGSERIAL NOT NULL PRIMARY KEY,
  
  game_id BIGINT NOT NULL DEFAULT 0,

  hokm BIGINT,

  hakem_index BIGINT NOT NULL DEFAULT 0,

  team_one_turn_score BIGINT Not NULL DEFAULT 0,
  team_two_turn_score BIGINT Not NULL DEFAULT 0,
  
  
  CONSTRAINT fk_game FOREIGN KEY (game_id) REFERENCES hokm4_game(id)
);

-- migrate:down

