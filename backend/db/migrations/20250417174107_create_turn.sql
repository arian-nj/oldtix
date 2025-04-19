-- migrate:up
CREATE TABLE turn (
  turn_id BIGSERIAL NOT NULL PRIMARY KEY,
  trick_id BIGINT NOT NULL DEFAULT 0,
  moves TEXT NOT NULL DEFAULT 0,
  CONSTRAINT fk_trick FOREIGN KEY (trick_id) REFERENCES trick(trick_id)
);

-- migrate:down

