-- name: GetPerson :one
SELECT * FROM person
WHERE id = $1 LIMIT 1;

-- name: GetPersonByUsername :one
SELECT * FROM person
WHERE username = $1 LIMIT 1;

-- name: InsertPerson :one
INSERT INTO person (
  username,display_name,hashed_password
) VALUES (
  $1,$2,$3
)
RETURNING *;

-- name: UpdatePersonDisplayName :exec
UPDATE person
SET display_name = $1
WHERE id = $2;

-- name: UpdateUserStatistics :exec
UPDATE user_statistic
SET 
    wins = wins + $1,
    losses = losses + $2,
    total_tricks_won = total_tricks_won + $3,
    total_tricks_lost = total_tricks_lost + $4,
    total_turns_won = total_turns_won + $5,
    total_turns_lost = total_turns_lost + $6,
    updated_at = NOW()
WHERE user_id = $7;