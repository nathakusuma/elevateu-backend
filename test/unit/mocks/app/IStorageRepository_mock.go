// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	io "io"

	mock "github.com/stretchr/testify/mock"
)

// MockIStorageRepository is an autogenerated mock type for the IStorageRepository type
type MockIStorageRepository struct {
	mock.Mock
}

type MockIStorageRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIStorageRepository) EXPECT() *MockIStorageRepository_Expecter {
	return &MockIStorageRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, path
func (_m *MockIStorageRepository) Delete(ctx context.Context, path string) error {
	ret := _m.Called(ctx, path)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIStorageRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockIStorageRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - path string
func (_e *MockIStorageRepository_Expecter) Delete(ctx interface{}, path interface{}) *MockIStorageRepository_Delete_Call {
	return &MockIStorageRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, path)}
}

func (_c *MockIStorageRepository_Delete_Call) Run(run func(ctx context.Context, path string)) *MockIStorageRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockIStorageRepository_Delete_Call) Return(_a0 error) *MockIStorageRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIStorageRepository_Delete_Call) RunAndReturn(run func(context.Context, string) error) *MockIStorageRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetSignedURL provides a mock function with given fields: path
func (_m *MockIStorageRepository) GetSignedURL(path string) (string, error) {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for GetSignedURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIStorageRepository_GetSignedURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSignedURL'
type MockIStorageRepository_GetSignedURL_Call struct {
	*mock.Call
}

// GetSignedURL is a helper method to define mock.On call
//   - path string
func (_e *MockIStorageRepository_Expecter) GetSignedURL(path interface{}) *MockIStorageRepository_GetSignedURL_Call {
	return &MockIStorageRepository_GetSignedURL_Call{Call: _e.mock.On("GetSignedURL", path)}
}

func (_c *MockIStorageRepository_GetSignedURL_Call) Run(run func(path string)) *MockIStorageRepository_GetSignedURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockIStorageRepository_GetSignedURL_Call) Return(_a0 string, _a1 error) *MockIStorageRepository_GetSignedURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIStorageRepository_GetSignedURL_Call) RunAndReturn(run func(string) (string, error)) *MockIStorageRepository_GetSignedURL_Call {
	_c.Call.Return(run)
	return _c
}

// Upload provides a mock function with given fields: ctx, file, path
func (_m *MockIStorageRepository) Upload(ctx context.Context, file io.Reader, path string) (string, error) {
	ret := _m.Called(ctx, file, path)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, io.Reader, string) (string, error)); ok {
		return rf(ctx, file, path)
	}
	if rf, ok := ret.Get(0).(func(context.Context, io.Reader, string) string); ok {
		r0 = rf(ctx, file, path)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, io.Reader, string) error); ok {
		r1 = rf(ctx, file, path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIStorageRepository_Upload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upload'
type MockIStorageRepository_Upload_Call struct {
	*mock.Call
}

// Upload is a helper method to define mock.On call
//   - ctx context.Context
//   - file io.Reader
//   - path string
func (_e *MockIStorageRepository_Expecter) Upload(ctx interface{}, file interface{}, path interface{}) *MockIStorageRepository_Upload_Call {
	return &MockIStorageRepository_Upload_Call{Call: _e.mock.On("Upload", ctx, file, path)}
}

func (_c *MockIStorageRepository_Upload_Call) Run(run func(ctx context.Context, file io.Reader, path string)) *MockIStorageRepository_Upload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(io.Reader), args[2].(string))
	})
	return _c
}

func (_c *MockIStorageRepository_Upload_Call) Return(_a0 string, _a1 error) *MockIStorageRepository_Upload_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIStorageRepository_Upload_Call) RunAndReturn(run func(context.Context, io.Reader, string) (string, error)) *MockIStorageRepository_Upload_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIStorageRepository creates a new instance of MockIStorageRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIStorageRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIStorageRepository {
	mock := &MockIStorageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
