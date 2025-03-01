-- name: GetPerson :one
SELECT * FROM person
WHERE id = $1 LIMIT 1;

-- name: GetPersonByUsername :one
SELECT * FROM person
WHERE username = $1 LIMIT 1;

-- name: InsertPerson :one
INSERT INTO person (
  username,hashed_password
) VALUES (
  $1,$2
)
RETURNING *;

-- name: UpdatePerson :exec
UPDATE person
SET display_name = $1,
bio = $2
WHERE id = $3;

