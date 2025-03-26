-- name: FindAll :many
SELECT * FROM posts;

-- name: FindByID :one
SELECT * FROM posts WHERE id = $1;

-- name: Create :one
INSERT INTO posts (author_id, title, content) VALUES ($1, $2, $3) RETURNING *;

-- name: Update :one
UPDATE posts SET title = $2, content = $3 WHERE id = $1 RETURNING *;

-- name: Delete :exec
DELETE FROM posts WHERE id = $1;