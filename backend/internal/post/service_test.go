package post_test

import (
	"backend/internal/post"
	"github.com/jackc/pgx/v5/pgtype"
	"testing"
	"time"
)

func TestService_GetPost(t *testing.T) {
	type args struct {
		id pgtype.UUID
	}

	test := []struct {
		name       string
		args       args
		setUp      func() (pgtype.UUID, post.Post)
		wantResult post.Post
		wantErr    bool
	}{
		{
			name: "Should return post",
			setUp: func() (pgtype.UUID, post.Post) {
				var id pgtype.UUID
				id.Scan("00000000-0000-0000-0000-000000000001")

				return id, post.Post{
					ID:       id,
					AuthorID: id,
					Title:    pgtype.Text{String: "Title"},
					Content:  pgtype.Text{String: "Content"},
					CreateAt: pgtype.Timestamptz{
						Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
