-- migrate:up
CREATE TABLE bot_user (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	tg_id TEXT NOT NULL UNIQUE,
	person_id BIGINT,
	CONSTRAINT fk_person FOREIGN KEY (person_id) REFERENCES person(id)
);

-- migrate:down

