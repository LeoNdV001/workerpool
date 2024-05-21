// Code generated by mockery. DO NOT EDIT.

package task

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// InterfaceTask is an autogenerated mock type for the InterfaceTask type
type InterfaceTask struct {
	mock.Mock
}

type InterfaceTask_Expecter struct {
	mock *mock.Mock
}

func (_m *InterfaceTask) EXPECT() *InterfaceTask_Expecter {
	return &InterfaceTask_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx
func (_m *InterfaceTask) Execute(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InterfaceTask_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type InterfaceTask_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
func (_e *InterfaceTask_Expecter) Execute(ctx interface{}) *InterfaceTask_Execute_Call {
	return &InterfaceTask_Execute_Call{Call: _e.mock.On("Execute", ctx)}
}

func (_c *InterfaceTask_Execute_Call) Run(run func(ctx context.Context)) *InterfaceTask_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *InterfaceTask_Execute_Call) Return(_a0 error) *InterfaceTask_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InterfaceTask_Execute_Call) RunAndReturn(run func(context.Context) error) *InterfaceTask_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// OnFailure provides a mock function with given fields: ctx, err
func (_m *InterfaceTask) OnFailure(ctx context.Context, err error) {
	_m.Called(ctx, err)
}

// InterfaceTask_OnFailure_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnFailure'
type InterfaceTask_OnFailure_Call struct {
	*mock.Call
}

// OnFailure is a helper method to define mock.On call
//   - ctx context.Context
//   - err error
func (_e *InterfaceTask_Expecter) OnFailure(ctx interface{}, err interface{}) *InterfaceTask_OnFailure_Call {
	return &InterfaceTask_OnFailure_Call{Call: _e.mock.On("OnFailure", ctx, err)}
}

func (_c *InterfaceTask_OnFailure_Call) Run(run func(ctx context.Context, err error)) *InterfaceTask_OnFailure_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(error))
	})
	return _c
}

func (_c *InterfaceTask_OnFailure_Call) Return() *InterfaceTask_OnFailure_Call {
	_c.Call.Return()
	return _c
}

func (_c *InterfaceTask_OnFailure_Call) RunAndReturn(run func(context.Context, error)) *InterfaceTask_OnFailure_Call {
	_c.Call.Return(run)
	return _c
}

// OnSuccess provides a mock function with given fields: ctx
func (_m *InterfaceTask) OnSuccess(ctx context.Context) {
	_m.Called(ctx)
}

// InterfaceTask_OnSuccess_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnSuccess'
type InterfaceTask_OnSuccess_Call struct {
	*mock.Call
}

// OnSuccess is a helper method to define mock.On call
//   - ctx context.Context
func (_e *InterfaceTask_Expecter) OnSuccess(ctx interface{}) *InterfaceTask_OnSuccess_Call {
	return &InterfaceTask_OnSuccess_Call{Call: _e.mock.On("OnSuccess", ctx)}
}

func (_c *InterfaceTask_OnSuccess_Call) Run(run func(ctx context.Context)) *InterfaceTask_OnSuccess_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *InterfaceTask_OnSuccess_Call) Return() *InterfaceTask_OnSuccess_Call {
	_c.Call.Return()
	return _c
}

func (_c *InterfaceTask_OnSuccess_Call) RunAndReturn(run func(context.Context)) *InterfaceTask_OnSuccess_Call {
	_c.Call.Return(run)
	return _c
}

// NewInterfaceTask creates a new instance of InterfaceTask. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInterfaceTask(t interface {
	mock.TestingT
	Cleanup(func())
}) *InterfaceTask {
	mock := &InterfaceTask{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
