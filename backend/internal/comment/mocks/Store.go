// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	comment "backend/internal/comment"
	context "context"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, arg
func (_m *Store) Create(ctx context.Context, arg comment.CreateRequest) (comment.Comment, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 comment.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, comment.CreateRequest) (comment.Comment, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, comment.CreateRequest) comment.Comment); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(comment.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, comment.CreateRequest) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx
func (_m *Store) GetAll(ctx context.Context) ([]comment.Comment, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []comment.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]comment.Comment, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []comment.Comment); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]comment.Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: ctx, id
func (_m *Store) GetById(ctx context.Context, id uuid.UUID) (comment.Comment, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetById")
	}

	var r0 comment.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (comment.Comment, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) comment.Comment); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(comment.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByPost provides a mock function with given fields: ctx, postId
func (_m *Store) GetByPost(ctx context.Context, postId uuid.UUID) ([]comment.Comment, error) {
	ret := _m.Called(ctx, postId)

	if len(ret) == 0 {
		panic("no return value specified for GetByPost")
	}

	var r0 []comment.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]comment.Comment, error)); ok {
		return rf(ctx, postId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []comment.Comment); ok {
		r0 = rf(ctx, postId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]comment.Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, postId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
