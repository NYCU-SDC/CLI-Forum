package post_test

import (
	"backend/internal/post"
	"backend/internal/post/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

//go:generate mockery --name=Querier --dir=./ --output=mocks
func TestService_GetPost(t *testing.T) {
	type args struct {
		id pgtype.UUID
	}

	test := []struct {
		name       string
		args       args
		setUp      func() (pgtype.UUID, pgtype.UUID, post.Post)
		wantResult post.Post
		wantErr    bool
	}{
		{
			name: "Should return post",
			setUp: func() (pgtype.UUID, pgtype.UUID, post.Post) {
				var id pgtype.UUID
				var authorID pgtype.UUID
				id.Scan("00000000-0000-0000-0000-000000000001")
				authorID.Scan("00000000-0000-0000-0000-000000000002")

				return id, authorID, post.Post{
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

	m := new(mocks.Querier)
	mockPost := post.Post{
		ID:       pgtype.UUID{},
		AuthorID: pgtype.UUID{},
		Title:    pgtype.Text{String: "Title"},
		Content:  pgtype.Text{String: "Content"},
		CreateAt: pgtype.Timestamptz{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	mockPost.ID.Scan("00000000-0000-0000-0000-000000000001")
	mockPost.AuthorID.Scan("00000000-0000-0000-0000-000000000001")
	m.On("GetPost", mock.Anything, mock.Anything).Return(, nil)

	for _, tt := range test {
		//t.Run(tt.name, func(t *testing.T) {
		//	id, authorI result := tt.setUp()
		//	tt.args.id = id
		//	tt.wantResult = result
		//
		//})
	}
}
