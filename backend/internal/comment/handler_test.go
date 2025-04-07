package comment_test

import (
	"backend/internal"
	"backend/internal/comment"
	"backend/internal/comment/mocks"
	"backend/internal/jwt"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_CreateHandler(t *testing.T) {
	type args struct {
		user          jwt.User
		requestPostId string
		requestBody   comment.CreateRequest
	}
	tests := []struct {
		name       string
		body       args
		wantStatus int
		wantResult comment.Response
	}{
		{
			name: "Should create comment",
			body: args{
				user: jwt.User{
					ID:       "7942c917-4770-43c1-a56a-952186b9970e",
					Username: "testuser",
					Role:     "user",
				},
				requestPostId: "7942c917-4770-43c1-a56a-952186b9970e",
				requestBody: comment.CreateRequest{
					Title:   "Test Title",
					Content: "Test Content",
				},
			},
			wantStatus: http.StatusOK,
			wantResult: comment.Response{
				ID:        "7942c917-4770-43c1-a56a-952186b9970e",
				PostId:    "7942c917-4770-43c1-a56a-952186b9970e",
				AuthorId:  "7942c917-4770-43c1-a56a-952186b9970e",
				Title:     "Test Title",
				Content:   "Test Content",
				CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC).String(),
			},
		},
		{
			name: "Should return error when title is empty",
			body: args{
				user: jwt.User{
					ID:       "7942c917-4770-43c1-a56a-952186b9970e",
					Username: "testuser",
					Role:     "user",
				},
				requestPostId: "7942c917-4770-43c1-a56a-952186b9970e",
				requestBody: comment.CreateRequest{
					Content: "Test Content",
				},
			},
			wantStatus: http.StatusBadRequest,
			wantResult: comment.Response{},
		},
		{
			name: "Should return error when content is empty",
			body: args{
				user: jwt.User{
					ID:       "7942c917-4770-43c1-a56a-952186b9970e",
					Username: "testuser",
					Role:     "user",
				},
				requestPostId: "7942c917-4770-43c1-a56a-952186b9970e",
				requestBody: comment.CreateRequest{
					Title: "Test Title",
				},
			},
			wantStatus: http.StatusBadRequest,
			wantResult: comment.Response{},
		},
		{
			name: "Should return error when postID is invalid",
			body: args{
				user: jwt.User{
					ID:       "7942c917-4770-43c1-a56a-952186b9970e",
					Username: "testuser",
					Role:     "user",
				},
				requestPostId: "7942c917-4770-43c1-952186b9970e",
				requestBody: comment.CreateRequest{
					Title: "Test Title",
				},
			},
			wantStatus: http.StatusBadRequest,
			wantResult: comment.Response{},
		},
	}

	// Mock the server
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("could not initialize logger: %v", err)
	}
	store := mocks.NewStore(t)
	store.On("Create", mock.Anything, comment.CreateRequest{
		PostID:   uuid.MustParse("7942c917-4770-43c1-a56a-952186b9970e"),
		AuthorID: uuid.MustParse("7942c917-4770-43c1-a56a-952186b9970e"),
		Title:    "Test Title",
		Content:  "Test Content",
	}).Return(
		comment.Comment{
			ID:        uuid.MustParse("7942c917-4770-43c1-a56a-952186b9970e"),
			PostID:    uuid.MustParse("7942c917-4770-43c1-a56a-952186b9970e"),
			AuthorID:  uuid.MustParse("7942c917-4770-43c1-a56a-952186b9970e"),
			Title:     pgtype.Text{String: "Test Title"},
			Content:   pgtype.Text{String: "Test Content"},
			CreatedAt: pgtype.Timestamptz{Time: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
		}, nil)

	h := comment.NewHandler(internal.NewValidator(), logger, store)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a requestBody
			requestBody, err := json.Marshal(tt.body.requestBody)
			if err != nil {
				t.Fatalf("could not marshal requestBody body: %v", err)
			}
			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/post/%s/comments", tt.body.requestPostId), bytes.NewReader(requestBody))
			r.SetPathValue("post_id", tt.body.requestPostId)
			w := httptest.NewRecorder()

			// Set the user in the context
			r = r.WithContext(context.WithValue(r.Context(), internal.UserContextKey, tt.body.user))

			// Call the handler
			h.CreateHandler(w, r)

			// Check the response
			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				res, err := json.Marshal(tt.wantResult)
				if err != nil {
					t.Fatalf("could not marshal want response: %v", err)
				}
				assert.Equal(t, string(res), strings.Trim(w.Body.String(), "\n"))
			}
		})
	}
}
