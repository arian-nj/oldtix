-- migrate:up
CREATE TABLE guest_person (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  user_id BIGINT Not NULL REFERENCES person(id),
  uid_string TEXT NOT NULL
);


-- migrate:down

