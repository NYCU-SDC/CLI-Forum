// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	post "backend/internal/post"
	context "context"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// Querier is an autogenerated mock type for the Querier type
type Querier struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, arg
func (_m *Querier) Create(ctx context.Context, arg post.CreateParams) (post.Post, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 post.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, post.CreateParams) (post.Post, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, post.CreateParams) post.Post); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(post.Post)
	}

	if rf, ok := ret.Get(1).(func(context.Context, post.CreateParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *Querier) Delete(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: ctx
func (_m *Querier) FindAll(ctx context.Context) ([]post.Post, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FindAll")
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

// FindByID provides a mock function with given fields: ctx, id
func (_m *Querier) FindByID(ctx context.Context, id uuid.UUID) (post.Post, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 post.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (post.Post, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) post.Post); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(post.Post)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, arg
func (_m *Querier) Update(ctx context.Context, arg post.UpdateParams) (post.Post, error) {
	ret := _m.Called(ctx, arg)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 post.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, post.UpdateParams) (post.Post, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, post.UpdateParams) post.Post); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Get(0).(post.Post)
	}

	if rf, ok := ret.Get(1).(func(context.Context, post.UpdateParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewQuerier creates a new instance of Querier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQuerier(t interface {
	mock.TestingT
	Cleanup(func())
}) *Querier {
	mock := &Querier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
