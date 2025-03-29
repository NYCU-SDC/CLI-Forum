-- name: FindAll :many
SELECT * FROM comments;

-- name: FindByID :one
SELECT * FROM comments WHERE id = $1;

-- name: Create :one
INSERT INTO comments (post_id, author_id, title, content) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: Update :one
UPDATE comments SET title = $2, content = $3 WHERE id = $1 RETURNING *;

-- name: Delete :exec
DELETE FROM comments WHERE id = $1;