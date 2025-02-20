// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/nathakusuma/elevateu-backend/domain/entity"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockIUserRepository is an autogenerated mock type for the IUserRepository type
type MockIUserRepository struct {
	mock.Mock
}

type MockIUserRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIUserRepository) EXPECT() *MockIUserRepository_Expecter {
	return &MockIUserRepository_Expecter{mock: &_m.Mock}
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *MockIUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIUserRepository_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type MockIUserRepository_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - user *entity.User
func (_e *MockIUserRepository_Expecter) CreateUser(ctx interface{}, user interface{}) *MockIUserRepository_CreateUser_Call {
	return &MockIUserRepository_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx, user)}
}

func (_c *MockIUserRepository_CreateUser_Call) Run(run func(ctx context.Context, user *entity.User)) *MockIUserRepository_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.User))
	})
	return _c
}

func (_c *MockIUserRepository_CreateUser_Call) Return(_a0 error) *MockIUserRepository_CreateUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIUserRepository_CreateUser_Call) RunAndReturn(run func(context.Context, *entity.User) error) *MockIUserRepository_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteUser provides a mock function with given fields: ctx, id
func (_m *MockIUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIUserRepository_DeleteUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUser'
type MockIUserRepository_DeleteUser_Call struct {
	*mock.Call
}

// DeleteUser is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *MockIUserRepository_Expecter) DeleteUser(ctx interface{}, id interface{}) *MockIUserRepository_DeleteUser_Call {
	return &MockIUserRepository_DeleteUser_Call{Call: _e.mock.On("DeleteUser", ctx, id)}
}

func (_c *MockIUserRepository_DeleteUser_Call) Run(run func(ctx context.Context, id uuid.UUID)) *MockIUserRepository_DeleteUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockIUserRepository_DeleteUser_Call) Return(_a0 error) *MockIUserRepository_DeleteUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIUserRepository_DeleteUser_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *MockIUserRepository_DeleteUser_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserByField provides a mock function with given fields: ctx, field, value
func (_m *MockIUserRepository) GetUserByField(ctx context.Context, field string, value string) (*entity.User, error) {
	ret := _m.Called(ctx, field, value)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByField")
	}

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*entity.User, error)); ok {
		return rf(ctx, field, value)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *entity.User); ok {
		r0 = rf(ctx, field, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, field, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIUserRepository_GetUserByField_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserByField'
type MockIUserRepository_GetUserByField_Call struct {
	*mock.Call
}

// GetUserByField is a helper method to define mock.On call
//   - ctx context.Context
//   - field string
//   - value string
func (_e *MockIUserRepository_Expecter) GetUserByField(ctx interface{}, field interface{}, value interface{}) *MockIUserRepository_GetUserByField_Call {
	return &MockIUserRepository_GetUserByField_Call{Call: _e.mock.On("GetUserByField", ctx, field, value)}
}

func (_c *MockIUserRepository_GetUserByField_Call) Run(run func(ctx context.Context, field string, value string)) *MockIUserRepository_GetUserByField_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockIUserRepository_GetUserByField_Call) Return(_a0 *entity.User, _a1 error) *MockIUserRepository_GetUserByField_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIUserRepository_GetUserByField_Call) RunAndReturn(run func(context.Context, string, string) (*entity.User, error)) *MockIUserRepository_GetUserByField_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUser provides a mock function with given fields: ctx, user
func (_m *MockIUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIUserRepository_UpdateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUser'
type MockIUserRepository_UpdateUser_Call struct {
	*mock.Call
}

// UpdateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - user *entity.User
func (_e *MockIUserRepository_Expecter) UpdateUser(ctx interface{}, user interface{}) *MockIUserRepository_UpdateUser_Call {
	return &MockIUserRepository_UpdateUser_Call{Call: _e.mock.On("UpdateUser", ctx, user)}
}

func (_c *MockIUserRepository_UpdateUser_Call) Run(run func(ctx context.Context, user *entity.User)) *MockIUserRepository_UpdateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.User))
	})
	return _c
}

func (_c *MockIUserRepository_UpdateUser_Call) Return(_a0 error) *MockIUserRepository_UpdateUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIUserRepository_UpdateUser_Call) RunAndReturn(run func(context.Context, *entity.User) error) *MockIUserRepository_UpdateUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIUserRepository creates a new instance of MockIUserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIUserRepository {
	mock := &MockIUserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
