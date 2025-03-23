package post

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"os"
)

func connect(ctx context.Context) (*pgx.Conn, error) {
	dbURL := os.Getenv("DATABASE_URL")
	return pgx.Connect(ctx, dbURL)
}

func GetAll(ctx context.Context) ([]Post, error) {
	// Connect to the database
	conn, err := connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	queries := New(conn)

	posts, err := queries.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func Get(ctx context.Context, id pgtype.UUID) (Post, error) {
	// Connect to the database
	conn, err := connect(ctx)
	if err != nil {
		return Post{}, err
	}
	defer conn.Close(ctx)

	queries := New(conn)

	return queries.FindByID(ctx, id)
}

func Create(ctx context.Context, post Post) (Post, error) {
	// Connect to the database
	conn, err := connect(ctx)
	if err != nil {
		return Post{}, err
	}
	defer conn.Close(ctx)

	queries := New(conn)

	return queries.Create(ctx, CreateParams{
		ID:       post.ID,
		AuthorID: post.AuthorID,
		Title:    post.Title,
		Content:  post.Content,
	})
}
