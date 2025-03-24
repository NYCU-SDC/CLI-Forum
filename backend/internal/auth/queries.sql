-- name: FindByName :one
SELECT * FROM users WHERE name = $1;

-- name: Create :one
INSERT INTO users (id, name, password) VALUES ($1, $2, $3) RETURNING *;

-- name: Update :one
UPDATE users SET name = $2, password = $3 WHERE id = $1 RETURNING *;

-- name: Delete :one
DELETE FROM users WHERE id = $1 RETURNING *;