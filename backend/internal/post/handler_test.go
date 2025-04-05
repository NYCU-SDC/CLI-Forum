package post_test

import (
	"backend/internal"
	"backend/internal/jwt"
	"backend/internal/post"
	"backend/internal/post/mocks"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_CreateHandler(t *testing.T) {
	type args struct {
		user    jwt.User
		request post.CreateRequest
	}

	tests := []struct {
		name       string
		args       args
		setupMock  func(m *mocks.Store)
		wantResult post.Response
		wantStatus int
	}{
		{
			name: "Should create post",
			args: args{
				user: jwt.User{
					ID:       "81c1ecc1-66d7-4134-b5fe-d886a6f418a4",
					Username: "test",
					Role:     "user",
				},
				request: post.CreateRequest{
					Title:   "Title",
					Content: "Content",
				},
			},
			setupMock: func(m *mocks.Store) {
				m.On("Create", mock.Anything, post.CreateRequest{
					AuthorID: uuid.MustParse("81c1ecc1-66d7-4134-b5fe-d886a6f418a4"),
					Title:    "Title",
					Content:  "Content",
				}).Return(post.Post{
					ID:       uuid.MustParse("54a46af2-b454-4746-8ab0-3cf26085a50b"),
					AuthorID: uuid.MustParse("81c1ecc1-66d7-4134-b5fe-d886a6f418a4"),
					Title:    pgtype.Text{String: "Title"},
					Content:  pgtype.Text{String: "Content"},
					CreateAt: pgtype.Timestamptz{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
				}, nil)
			},
			wantResult: post.Response{
				ID:       "54a46af2-b454-4746-8ab0-3cf26085a50b",
				AuthorID: "81c1ecc1-66d7-4134-b5fe-d886a6f418a4",
				Title:    "Title",
				Content:  "Content",
				CreateAt: "2000-01-01T00:00:00Z",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Should return error when content is empty",
			args: args{
				user: jwt.User{
					ID:       "81c1ecc1-66d7-4134-b5fe-d886a6f418a4",
					Username: "test",
					Role:     "user",
				},
				request: post.CreateRequest{
					Title:   "Title",
					Content: "",
				},
			},
			setupMock: func(m *mocks.Store) {
				m.On("Create", mock.Anything, post.CreateRequest{
					AuthorID: uuid.MustParse("81c1ecc1-66d7-4134-b5fe-d886a6f418a4"),
					Title:    "Title",
					Content:  "Content",
				}).Return(post.Post{
					ID:       uuid.MustParse("54a46af2-b454-4746-8ab0-3cf26085a50b"),
					AuthorID: uuid.MustParse("81c1ecc1-66d7-4134-b5fe-d886a6f418a4"),
					Title:    pgtype.Text{String: "Title"},
					Content:  pgtype.Text{String: "Content"},
					CreateAt: pgtype.Timestamptz{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
				}, nil)
			},
			wantResult: post.Response{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Should return error when title is empty",
			args: args{
				user: jwt.User{
					ID:       "81c1ecc1-66d7-4134-b5fe-d886a6f418a4",
					Username: "test",
					Role:     "user",
				},
				request: post.CreateRequest{
					Title:   "",
					Content: "Content",
				},
			},
			setupMock: func(m *mocks.Store) {
				m.On("Create", mock.Anything, post.CreateRequest{
					AuthorID: uuid.MustParse("81c1ecc1-66d7-4134-b5fe-d886a6f418a4"),
					Title:    "Title",
					Content:  "Content",
				}).Return(post.Post{
					ID:       uuid.MustParse("54a46af2-b454-4746-8ab0-3cf26085a50b"),
					AuthorID: uuid.MustParse("81c1ecc1-66d7-4134-b5fe-d886a6f418a4"),
					Title:    pgtype.Text{String: "Title"},
					Content:  pgtype.Text{String: "Content"},
					CreateAt: pgtype.Timestamptz{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
				}, nil)
			},
			wantResult: post.Response{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mocks.Store)
			tt.setupMock(m)

			requestBody, err := json.Marshal(tt.args.request)
			if err != nil {
				assert.Failf(t, "Failed to marshal request body", "%+v", err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/posts", bytes.NewReader(requestBody))
			ctx := context.WithValue(r.Context(), internal.UserContextKey, tt.args.user)
			r = r.WithContext(ctx)

			logger, err := zap.NewDevelopment()
			if err != nil {
				assert.Failf(t, "Failed to marshal request body", "%+v", err)
			}
			h := post.NewHandler(validator.New(), logger, m)

			h.CreateHandler(w, r)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusBadRequest {
				assert.Contains(t, w.Body.String(), "Validation Problem")
			} else {
				jsonWant, err := json.Marshal(tt.wantResult)
				if err != nil {
					t.Fatalf("failed to marshal expected response: %v", err)
				}
				assert.Equal(t, string(jsonWant), w.Body.String())
			}
		})
	}
}

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
