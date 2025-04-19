-- migrate:up
CREATE TABLE user_statistic(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT Not NULL REFERENCES person(id),
    win BIGINT Not NULL DEFAULT 0,
    lose BIGINT Not NULL DEFAULT 0,
    total_tricks_won BIGINT Not NULL DEFAULT 0,
    total_tricks_lost BIGINT Not NULL DEFAULT 0,
    total_turns_won BIGINT Not NULL DEFAULT 0,
    total_turns_lost BIGINT Not NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW()
);


-- migrate:down

