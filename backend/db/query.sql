-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: InsertUser :one
INSERT INTO users (
  username,hashed_password
) VALUES (
  $1,$2
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users 
SET display_name = $1,
bio = $2
WHERE id = $3;

-- -- name: ListAuthors :many
-- SELECT * FROM authors
-- ORDER BY name;


-- -- name: UpdateAuthor :exec
-- UPDATE authors
--   set name = $2,
--   bio = $3
-- WHERE id = $1;

-- -- name: DeleteAuthor :exec
-- DELETE FROM authors
-- WHERE id = $1;