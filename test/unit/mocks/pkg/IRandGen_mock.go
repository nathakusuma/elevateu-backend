// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MockIRandGen is an autogenerated mock type for the IRandGen type
type MockIRandGen struct {
	mock.Mock
}

type MockIRandGen_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIRandGen) EXPECT() *MockIRandGen_Expecter {
	return &MockIRandGen_Expecter{mock: &_m.Mock}
}

// RandomNumber provides a mock function with given fields: digits
func (_m *MockIRandGen) RandomNumber(digits int) (int, error) {
	ret := _m.Called(digits)

	if len(ret) == 0 {
		panic("no return value specified for RandomNumber")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (int, error)); ok {
		return rf(digits)
	}
	if rf, ok := ret.Get(0).(func(int) int); ok {
		r0 = rf(digits)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(digits)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRandGen_RandomNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RandomNumber'
type MockIRandGen_RandomNumber_Call struct {
	*mock.Call
}

// RandomNumber is a helper method to define mock.On call
//   - digits int
func (_e *MockIRandGen_Expecter) RandomNumber(digits interface{}) *MockIRandGen_RandomNumber_Call {
	return &MockIRandGen_RandomNumber_Call{Call: _e.mock.On("RandomNumber", digits)}
}

func (_c *MockIRandGen_RandomNumber_Call) Run(run func(digits int)) *MockIRandGen_RandomNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockIRandGen_RandomNumber_Call) Return(_a0 int, _a1 error) *MockIRandGen_RandomNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRandGen_RandomNumber_Call) RunAndReturn(run func(int) (int, error)) *MockIRandGen_RandomNumber_Call {
	_c.Call.Return(run)
	return _c
}

// RandomString provides a mock function with given fields: length
func (_m *MockIRandGen) RandomString(length int) (string, error) {
	ret := _m.Called(length)

	if len(ret) == 0 {
		panic("no return value specified for RandomString")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (string, error)); ok {
		return rf(length)
	}
	if rf, ok := ret.Get(0).(func(int) string); ok {
		r0 = rf(length)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(length)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRandGen_RandomString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RandomString'
type MockIRandGen_RandomString_Call struct {
	*mock.Call
}

// RandomString is a helper method to define mock.On call
//   - length int
func (_e *MockIRandGen_Expecter) RandomString(length interface{}) *MockIRandGen_RandomString_Call {
	return &MockIRandGen_RandomString_Call{Call: _e.mock.On("RandomString", length)}
}

func (_c *MockIRandGen_RandomString_Call) Run(run func(length int)) *MockIRandGen_RandomString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockIRandGen_RandomString_Call) Return(_a0 string, _a1 error) *MockIRandGen_RandomString_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRandGen_RandomString_Call) RunAndReturn(run func(int) (string, error)) *MockIRandGen_RandomString_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIRandGen creates a new instance of MockIRandGen. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIRandGen(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIRandGen {
	mock := &MockIRandGen{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
