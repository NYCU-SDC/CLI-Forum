package post_test

import (
	"backend/internal/post"
	"backend/internal/post/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetPost(t *testing.T) {
	type args struct {
		id string
	}

	tests := []struct {
		name             string
		args             args
		wantResult       *post.Response
		wantResponseCode int
	}{

		{
			name:             "Should return post",
			args:             args{},
			wantResult:       &post.Response{},
			wantResponseCode: http.StatusOK,
		},
		{
			name: "Should return error when id is could not be parsed",
			args: args{
				id: "1",
			},
			wantResult:       nil,
			wantResponseCode: http.StatusBadRequest,
		},
	}

	m := new(mocks.Servicer)
	m.On("GetPost", mock.Anything, mock.Anything).Return(post.Post{}, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := post.NewHandler(m, zap.NewExample().Sugar())
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/post?id="+tt.args.id, nil)
			h.GetPost(w, r)
			assert.Equal(t, tt.wantResponseCode, w.Code)
		})
	}
}
