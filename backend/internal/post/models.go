// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package post

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Post struct {
	ID       pgtype.UUID
	AuthorID pgtype.UUID
	Title    pgtype.Text
	Content  pgtype.Text
	CreateAt pgtype.Timestamptz
}
