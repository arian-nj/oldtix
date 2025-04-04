-- name: InsertHokm4Game :one
INSERT INTO hokm4_game DEFAULT VALUES RETURNING *;

-- name: InsertGamePlayer :one
INSERT INTO game_player (player_id, game_id,team) VALUES ($1,$2,$3) RETURNING *;

-- name: InsertTrick :one
INSERT INTO trick (game_id, hakem_index) VALUES ($1,$2) RETURNING *;

-- name: UpdateHokmTrick :exec
UPDATE trick SET hokm = $1 WHERE trick_id = $2;

-- name: InsertTurn :one
INSERT INTO turn (moves,trick_id) VALUES ($1,$2) RETURNING *;


-- name: UpdateTrickScores :exec
UPDATE hokm4_game SET team_one_tricks_score = $1,team_two_tricks_score=$2 WHERE id = $3;

-- name: UpdateTurnScores :exec
UPDATE trick SET team_one_turn_score = $1,team_two_turn_score=$2 WHERE trick_id = $3;

-- -- name: InsertHokm4Statistic :exec
-- INSERT INTO hokm4_game_statistic (match_id,person_id,tricks_won,tricks_lost,turns_won,turns_lost,is_won) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING *;