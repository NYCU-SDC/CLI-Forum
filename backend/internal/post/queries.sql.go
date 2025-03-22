// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package post

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const create = `-- name: Create :one
INSERT INTO posts (id, author_id, title, content) VALUES ($1, $2, $3, $4) RETURNING id, author_id, title, content, create_at
`

type CreateParams struct {
	ID       pgtype.UUID
	AuthorID pgtype.UUID
	Title    pgtype.Text
	Content  pgtype.Text
}

func (q *Queries) Create(ctx context.Context, arg CreateParams) (Post, error) {
	row := q.db.QueryRow(ctx, create,
		arg.ID,
		arg.AuthorID,
		arg.Title,
		arg.Content,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Title,
		&i.Content,
		&i.CreateAt,
	)
	return i, err
}

const delete = `-- name: Delete :exec
DELETE FROM posts WHERE id = $1
`

func (q *Queries) Delete(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, delete, id)
	return err
}

const findAll = `-- name: FindAll :many
SELECT id, author_id, title, content, create_at FROM posts
`

func (q *Queries) FindAll(ctx context.Context) ([]Post, error) {
	rows, err := q.db.Query(ctx, findAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.AuthorID,
			&i.Title,
			&i.Content,
			&i.CreateAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findByID = `-- name: FindByID :one
SELECT id, author_id, title, content, create_at FROM posts WHERE id = $1
`

func (q *Queries) FindByID(ctx context.Context, id pgtype.UUID) (Post, error) {
	row := q.db.QueryRow(ctx, findByID, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Title,
		&i.Content,
		&i.CreateAt,
	)
	return i, err
}

const update = `-- name: Update :one
UPDATE posts SET title = $2, content = $3 WHERE id = $1 RETURNING id, author_id, title, content, create_at
`

type UpdateParams struct {
	ID      pgtype.UUID
	Title   pgtype.Text
	Content pgtype.Text
}

func (q *Queries) Update(ctx context.Context, arg UpdateParams) (Post, error) {
	row := q.db.QueryRow(ctx, update, arg.ID, arg.Title, arg.Content)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Title,
		&i.Content,
		&i.CreateAt,
	)
	return i, err
}
