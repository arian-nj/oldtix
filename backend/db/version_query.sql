-- name: GetVersion :one
SELECT * FROM patch_version WHERE rmode = $1 LIMIT 1;

-- name: InsertVersion :one
INSERT INTO patch_version (version_number,rmode) VALUES ($1,$2)
RETURNING *;

-- name: UpdateVersion :one
UPDATE patch_version SET version_number = $1 WHERE rmode = $2 
RETURNING *;