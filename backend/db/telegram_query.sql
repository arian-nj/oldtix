-- name: InsertBotUser :one
INSERT INTO bot_user (
	tg_id
) VALUES (
	$1
) ON CONFLICT DO NOTHING
RETURNING *;

