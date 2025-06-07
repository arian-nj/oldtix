-- name: GetPerson :one
SELECT * FROM person
WHERE id = $1 LIMIT 1;

-- -- name: GetPersonByUsername :one
-- SELECT * FROM person
-- WHERE username = $1 LIMIT 1;

-- name: InsertPerson :one
INSERT INTO person (
  display_name,coin
) VALUES (
  $1,$2
) RETURNING *;

-- name: UpdatePersonDisplayName :exec
UPDATE person
SET display_name = $1
WHERE id = $2;

-- name: AddCoinToPerson :exec
UPDATE person
SET coin = coin + $1
WHERE id = $2;

-- name: InsertUserStatistic :one
INSERT INTO user_statistic (
  user_id
) VALUES (
  $1
) 
RETURNING *;

-- name: UpdateUserStatistics :exec
UPDATE user_statistic
SET 
    win = win + $1,
    lose = lose + $2,
    total_tricks_won = total_tricks_won + $3,
    total_tricks_lost = total_tricks_lost + $4,
    total_turns_won = total_turns_won + $5,
    total_turns_lost = total_turns_lost + $6,
    updated_at = NOW()
WHERE user_id = $7;

-- name: GetPersonStatisticsById :one
SELECT * FROM user_statistic WHERE user_id = $1;

-- name: InsertGuestPerson :one
INSERT INTO guest_person (
    uid_string,user_id
  ) VALUES (
    $1,$2
  ) RETURNING *;

-- name: GetGuestPersonByUid :one
SELECT * FROM guest_person WHERE uid_string = $1;
