package comment_test

import (
	"backend/internal/comment"
	"backend/internal/comment/mocks"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func MustParseUUID(t *testing.T, src string) pgtype.UUID {
	var uuid pgtype.UUID
	err := uuid.Scan(src)
	if err != nil {
		t.Fatalf("Invalid UUID %s", err.Error())
	}
	return uuid
}

func MustParseTime(t *testing.T, src time.Time) pgtype.Timestamptz {
	var timestamptz pgtype.Timestamptz
	err := timestamptz.Scan(src)
	if err != nil {
		t.Fatalf("Invalid time %s", err.Error())
	}
	return timestamptz
}

func TestService_Create(t *testing.T) {
	// create a test case
	testCases := []struct {
		name      string
		data      comment.CreateParams
		expectRes comment.Comment
		expectErr error
	}{
		{"Correct data with all fields", comment.CreateParams{
			PostID:   MustParseUUID(t, "00000000-0000-0000-0000-000000000001"),
			AuthorID: MustParseUUID(t, "00000000-0000-0000-0000-000000000002"),
			Title:    "Test Title",
			Content:  pgtype.Text{String: "Test Content", Valid: true},
		}, comment.Comment{
			ID:        MustParseUUID(t, "00000000-0000-0000-0000-000000000003"),
			PostID:    MustParseUUID(t, "00000000-0000-0000-0000-000000000001"),
			AuthorID:  MustParseUUID(t, "00000000-0000-0000-0000-000000000002"),
			Title:     "Test Title",
			Content:   pgtype.Text{String: "Test Content", Valid: true},
			CreatedAt: MustParseTime(t, time.Now()),
		}, nil},
		{"Error occur when creating comment", comment.CreateParams{
			PostID:   MustParseUUID(t, "00000000-0000-0000-0000-000000000111"),
			AuthorID: MustParseUUID(t, "00000000-0000-0000-0000-000000000112"),
		},
			comment.Comment{}, errors.New("error creating comment")},
	}

	// create a service with a mock querier
	querier := mocks.NewQuerier(t)
	service := comment.NewService(zap.NewExample(), querier)

	// set up the mock expectation
	querier.On("Create", context.Background(), testCases[0].data).Return(comment.Comment{
		ID:        MustParseUUID(t, "00000000-0000-0000-0000-000000000003"),
		PostID:    testCases[0].data.PostID,
		AuthorID:  testCases[0].data.AuthorID,
		Title:     testCases[0].data.Title,
		Content:   testCases[0].data.Content,
		CreatedAt: testCases[0].expectRes.CreatedAt,
	}, nil).Once()
	querier.On("Create", context.Background(), testCases[1].data).Return(comment.Comment{}, errors.New("error creating comment")).Once()

	// iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := service.Create(context.Background(), tc.data)
			assert.Equal(t, tc.expectRes, res)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}
