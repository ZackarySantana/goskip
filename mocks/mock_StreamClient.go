// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	skip "github.com/zackarysantana/goskip"
)

// MockStreamClient is an autogenerated mock type for the StreamClient type
type MockStreamClient struct {
	mock.Mock
}

type MockStreamClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockStreamClient) EXPECT() *MockStreamClient_Expecter {
	return &MockStreamClient_Expecter{mock: &_m.Mock}
}

// Stream provides a mock function with given fields: ctx, uuid, callback
func (_m *MockStreamClient) Stream(ctx context.Context, uuid string, callback func(skip.StreamType, []byte) error) error {
	ret := _m.Called(ctx, uuid, callback)

	if len(ret) == 0 {
		panic("no return value specified for Stream")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, func(skip.StreamType, []byte) error) error); ok {
		r0 = rf(ctx, uuid, callback)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockStreamClient_Stream_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stream'
type MockStreamClient_Stream_Call struct {
	*mock.Call
}

// Stream is a helper method to define mock.On call
//   - ctx context.Context
//   - uuid string
//   - callback func(skip.StreamType , []byte) error
func (_e *MockStreamClient_Expecter) Stream(ctx interface{}, uuid interface{}, callback interface{}) *MockStreamClient_Stream_Call {
	return &MockStreamClient_Stream_Call{Call: _e.mock.On("Stream", ctx, uuid, callback)}
}

func (_c *MockStreamClient_Stream_Call) Run(run func(ctx context.Context, uuid string, callback func(skip.StreamType, []byte) error)) *MockStreamClient_Stream_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(func(skip.StreamType, []byte) error))
	})
	return _c
}

func (_c *MockStreamClient_Stream_Call) Return(_a0 error) *MockStreamClient_Stream_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockStreamClient_Stream_Call) RunAndReturn(run func(context.Context, string, func(skip.StreamType, []byte) error) error) *MockStreamClient_Stream_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockStreamClient creates a new instance of MockStreamClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockStreamClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockStreamClient {
	mock := &MockStreamClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
