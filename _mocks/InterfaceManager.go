// Code generated by mockery. DO NOT EDIT.

package workerpool

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	"github.com/LeoNdV001/workerpool"
)

// InterfaceManager is an autogenerated mock type for the InterfaceManager type
type InterfaceManager struct {
	mock.Mock
}

type InterfaceManager_Expecter struct {
	mock *mock.Mock
}

func (_m *InterfaceManager) EXPECT() *InterfaceManager_Expecter {
	return &InterfaceManager_Expecter{mock: &_m.Mock}
}

// NewWorkerPool provides a mock function with given fields: ctx
func (_m *InterfaceManager) NewWorkerPool(ctx context.Context) workerpool.InterfaceWorkerPool {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for NewWorkerPool")
	}

	var r0 workerpool.InterfaceWorkerPool
	if rf, ok := ret.Get(0).(func(context.Context) workerpool.InterfaceWorkerPool); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(workerpool.InterfaceWorkerPool)
		}
	}

	return r0
}

// InterfaceManager_NewWorkerPool_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewWorkerPool'
type InterfaceManager_NewWorkerPool_Call struct {
	*mock.Call
}

// NewWorkerPool is a helper method to define mock.On call
//   - ctx context.Context
func (_e *InterfaceManager_Expecter) NewWorkerPool(ctx interface{}) *InterfaceManager_NewWorkerPool_Call {
	return &InterfaceManager_NewWorkerPool_Call{Call: _e.mock.On("NewWorkerPool", ctx)}
}

func (_c *InterfaceManager_NewWorkerPool_Call) Run(run func(ctx context.Context)) *InterfaceManager_NewWorkerPool_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *InterfaceManager_NewWorkerPool_Call) Return(_a0 workerpool.InterfaceWorkerPool) *InterfaceManager_NewWorkerPool_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InterfaceManager_NewWorkerPool_Call) RunAndReturn(run func(context.Context) workerpool.InterfaceWorkerPool) *InterfaceManager_NewWorkerPool_Call {
	_c.Call.Return(run)
	return _c
}

// NewInterfaceManager creates a new instance of InterfaceManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInterfaceManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *InterfaceManager {
	mock := &InterfaceManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
