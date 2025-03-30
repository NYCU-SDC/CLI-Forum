-- name: GetByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetByName :one
SELECT * FROM users WHERE name = $1;

-- name: Create :one
INSERT INTO users (name, password) VALUES ($1, $2) RETURNING *;

-- name: UpdateName :one
UPDATE users SET name = $2, password = $3 WHERE id = $1 RETURNING *;

-- name: UpdatePassword :execrows
UPDATE users SET password = $2 WHERE id = $1;

-- name: Delete :execrows
DELETE FROM users WHERE id = $1;