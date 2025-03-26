// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	pgtype "github.com/jackc/pgx/v5/pgtype"
	mock "github.com/stretchr/testify/mock"

	post "backend/internal/post"
)

// Finder is an autogenerated mock type for the Finder type
type Finder struct {
	mock.Mock
}

// CreatePost provides a mock function with given fields: ctx, request
func (_m *Finder) CreatePost(ctx context.Context, request post.CreateRequest) (post.Post, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for CreatePost")
	}

	var r0 post.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, post.CreateRequest) (post.Post, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, post.CreateRequest) post.Post); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(post.Post)
	}

	if rf, ok := ret.Get(1).(func(context.Context, post.CreateRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx
func (_m *Finder) GetAll(ctx context.Context) ([]post.Post, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []post.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]post.Post, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []post.Post); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]post.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPost provides a mock function with given fields: ctx, id
func (_m *Finder) GetPost(ctx context.Context, id pgtype.UUID) (post.Post, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetPost")
	}

	var r0 post.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, pgtype.UUID) (post.Post, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pgtype.UUID) post.Post); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(post.Post)
	}

	if rf, ok := ret.Get(1).(func(context.Context, pgtype.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFinder creates a new instance of Finder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFinder(t interface {
	mock.TestingT
	Cleanup(func())
}) *Finder {
	mock := &Finder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
