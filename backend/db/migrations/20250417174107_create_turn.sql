-- migrate:up
CREATE TABLE turn (
  turn_id BIGSERIAL NOT NULL PRIMARY KEY,
  trick_id BIGINT NOT NULL,
  moves TEXT NOT NULL,
  CONSTRAINT fk_trick FOREIGN KEY (trick_id) REFERENCES trick(trick_id)
);

-- migrate:down

