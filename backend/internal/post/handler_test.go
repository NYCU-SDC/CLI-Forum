package post_test

import (
	"backend/internal/post"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateResponse(t *testing.T) {
	tests := []struct {
		name       string
		post       post.Post
		wantResult post.Response
	}{
		{
			name: "Should return post",
			post: post.Post{
				ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				AuthorID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Title:    pgtype.Text{String: "Title"},
				Content:  pgtype.Text{String: "Content"},
				CreateAt: pgtype.Timestamptz{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			wantResult: post.Response{
				ID:       "00000000-0000-0000-0000-000000000001",
				AuthorID: "00000000-0000-0000-0000-000000000001",
				Title:    "Title",
				Content:  "Content",
				CreateAt: "2000-01-01T00:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := post.GenerateResponse(tt.post)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}
