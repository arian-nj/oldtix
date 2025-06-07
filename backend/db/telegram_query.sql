-- name: InsertBotUser :one
INSERT INTO bot_user (
	tg_id
) VALUES (
	$1
) ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetBotUser :one
SELECT * FROM bot_user WHERE tg_id = $1;

-- name: UpdateBotUsersPersonId :exec
UPDATE bot_user 
SET person_id = $1 WHERE tg_id = $2;
