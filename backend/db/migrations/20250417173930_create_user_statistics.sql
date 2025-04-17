-- migrate:up
CREATE TABLE user_statistic(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT Not NULL REFERENCES person(id),
    win BIGINT Not NULL,
    lose BIGINT Not NULL,
    total_tricks_won BIGINT Not NULL,
    total_tricks_lost BIGINT Not NULL,
    total_turns_won BIGINT Not NULL,
    total_turns_lost BIGINT Not NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);


-- migrate:down

