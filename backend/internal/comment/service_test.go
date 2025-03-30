package comment_test

import (
	"backend/internal/comment"
	"backend/internal/comment/mocks"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	fixedTime := time.Date(2023, 10, 10, 10, 0, 0, 0, time.UTC)
	errCreate := errors.New("error creating comment")

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
			CreatedAt: MustParseTime(t, fixedTime),
		}, nil},
		{"Error occur when creating comment", comment.CreateParams{
			PostID:   MustParseUUID(t, "00000000-0000-0000-0000-000000000111"),
			AuthorID: MustParseUUID(t, "00000000-0000-0000-0000-000000000112"),
		},
			comment.Comment{}, errCreate},
	}

	// create a service with a mock querier
	querier := mocks.NewQuerier(t)
	service := comment.NewService(zap.NewExample(), querier)

	// set up the mock expectation
	querier.On("Create", mock.Anything, testCases[0].data).Return(testCases[0].expectRes, nil).Once()
	querier.On("Create", mock.Anything, testCases[1].data).Return(comment.Comment{}, errCreate).Once()

	// iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := service.Create(context.Background(), tc.data)
			assert.Equal(t, tc.expectRes, res)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestService_GetAll(t *testing.T) {
	// create test cases
	fixedTime := time.Date(2023, 10, 10, 10, 0, 0, 0, time.UTC)
	errGetAll := errors.New("error getting all comments")

	testCases := []struct {
		name      string
		ctx       context.Context
		expectRes []comment.Comment
		expectErr error
	}{
		{"Correct get all comments", context.Background(), []comment.Comment{
			{
				ID:        MustParseUUID(t, "00000000-0000-0000-0000-000000000001"),
				PostID:    MustParseUUID(t, "00000000-0000-0000-0000-000000000002"),
				AuthorID:  MustParseUUID(t, "00000000-0000-0000-0000-000000000003"),
				Title:     "Test Title 1",
				Content:   pgtype.Text{String: "Test Content 1", Valid: true},
				CreatedAt: MustParseTime(t, fixedTime),
			},
			{
				ID:        MustParseUUID(t, "00000000-0000-0000-0000-000000000002"),
				PostID:    MustParseUUID(t, "00000000-0000-0000-0000-000000000003"),
				AuthorID:  MustParseUUID(t, "00000000-0000-0000-0000-000000000004"),
				Title:     "Test Title 2",
				Content:   pgtype.Text{String: "Test Content 2", Valid: true},
				CreatedAt: MustParseTime(t, fixedTime),
			},
		}, nil,
		},
		{"Error occur when getting all comments", context.Background(), nil, errGetAll},
	}

	// create a service with a mock querier
	querier := mocks.NewQuerier(t)
	service := comment.NewService(zap.NewExample(), querier)

	// set up the mock expectation
	querier.On("FindAll", testCases[0].ctx).Return(testCases[0].expectRes, nil).Once()
	querier.On("FindAll", testCases[1].ctx).Return([]comment.Comment{}, errGetAll).Once()

	// iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := service.GetAll(tc.ctx)
			assert.Equal(t, tc.expectRes, res)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestService_GetById(t *testing.T) {
	// create test cases
	fixedTime := time.Date(2023, 10, 10, 10, 0, 0, 0, time.UTC)
	errGetByID := errors.New("error getting comment by ID")

	testCases := []struct {
		name      string
		id        pgtype.UUID
		expectRes comment.Comment
		expectErr error
	}{
		{
			"Correct get comment by ID",
			MustParseUUID(t, "00000000-0000-0000-0000-000000000001"),
			comment.Comment{
				ID:        MustParseUUID(t, "00000000-0000-0000-0000-000000000001"),
				PostID:    MustParseUUID(t, "00000000-0000-0000-0000-000000000002"),
				AuthorID:  MustParseUUID(t, "00000000-0000-0000-0000-000000000003"),
				Title:     "Test Title 1",
				Content:   pgtype.Text{String: "Test Content 1", Valid: true},
				CreatedAt: MustParseTime(t, fixedTime),
			}, nil,
		},
		{
			"Error getting comment by ID",
			MustParseUUID(t, "00000000-0000-0000-0000-000000000002"),
			comment.Comment{}, errGetByID,
		},
	}

	// create a service with a mock querier
	querier := mocks.NewQuerier(t)
	service := comment.NewService(zap.NewExample(), querier)

	// set up the mock expectation
	querier.On("FindByID", mock.Anything, testCases[0].id).Return(testCases[0].expectRes, nil).Once()
	querier.On("FindByID", mock.Anything, testCases[1].id).Return(testCases[1].expectRes, errGetByID).Once()

	// iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := service.GetById(context.Background(), tc.id)
			assert.Equal(t, tc.expectRes, res)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestService_Delete(t *testing.T) {
	// create test cases
	errDelete := errors.New("error getting comment by ID")

	testCases := []struct {
		name      string
		id        pgtype.UUID
		expectErr error
	}{
		{
			"Correct delete comment",
			MustParseUUID(t, "00000000-0000-0000-0000-000000000001"),
			nil,
		},
		{
			"Error getting comment by ID",
			MustParseUUID(t, "00000000-0000-0000-0000-000000000002"),
			errDelete,
		},
	}

	// create a service with a mock querier
	querier := mocks.NewQuerier(t)
	service := comment.NewService(zap.NewExample(), querier)

	// set up the mock expectation
	querier.On("Delete", mock.Anything, testCases[0].id).Return(nil).Once()
	querier.On("Delete", mock.Anything, testCases[1].id).Return(errDelete).Once()

	// iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.Delete(context.Background(), tc.id)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestService_Update(t *testing.T) {
	// create test cases
	fixedTime := time.Date(2023, 10, 10, 10, 0, 0, 0, time.UTC)
	errUpdate := errors.New("error updating comment")

	testCases := []struct {
		name      string
		data      comment.UpdateParams
		expectRes comment.Comment
		expectErr error
	}{
		{
			"Correct update comment",
			comment.UpdateParams{
				ID:      MustParseUUID(t, "00000000-0000-0000-0000-000000000001"),
				Title:   "Updated Title",
				Content: pgtype.Text{String: "Updated Content", Valid: true},
			},
			comment.Comment{
				ID:        MustParseUUID(t, "00000000-0000-0000-0000-000000000001"),
				PostID:    MustParseUUID(t, "00000000-0000-0000-0000-000000000002"),
				AuthorID:  MustParseUUID(t, "00000000-0000-0000-0000-000000000003"),
				Title:     "Updated Title",
				Content:   pgtype.Text{String: "Updated Content", Valid: true},
				CreatedAt: MustParseTime(t, fixedTime),
			}, nil,
		},
		{
			"Error updating comment",
			comment.UpdateParams{
				ID:      MustParseUUID(t, "00000000-0000-0000-0000-000000000001"),
				Title:   "Updated Title",
				Content: pgtype.Text{String: "Updated Content", Valid: true},
			},
			comment.Comment{}, errUpdate,
		},
	}

	// create a service with a mock querier
	querier := mocks.NewQuerier(t)
	service := comment.NewService(zap.NewExample(), querier)

	// set up the mock expectation
	querier.On("Update", mock.Anything, testCases[0].data).Return(testCases[0].expectRes, nil).Once()
	querier.On("Update", mock.Anything, testCases[1].data).Return(testCases[1].expectRes, errUpdate).Once()

	// iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := service.Update(context.Background(), tc.data)
			assert.Equal(t, tc.expectRes, res)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}
