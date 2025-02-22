-- -- name: GetVersion :one
-- SELECT * FROM patch_version LIMIT 1;

-- -- name: InsertVersion :one
-- INSERT INTO patch_version (version_number) VALUES ($1)
-- RETURNING *;

-- -- name: UpdateVersion :exec
-- UPDATE patch_version SET version_number = $1 WHERE id = $2;