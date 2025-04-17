-- migrate:up
CREATE TABLE trick (
  trick_id BIGSERIAL NOT NULL PRIMARY KEY,
  
  game_id BIGINT NOT NULL,

  hokm BIGINT,

  hakem_index BIGINT NOT NULL,

  team_one_turn_score BIGINT Not NULL,
  team_two_turn_score BIGINT Not NULL,
  
  CONSTRAINT fk_game FOREIGN KEY (game_id) REFERENCES hokm4_game(id)
);

-- migrate:down

