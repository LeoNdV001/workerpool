package examples

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/LeoNdV001/workerpool"
	mocks "github.com/LeoNdV001/workerpool/_mocks"
)

// SimpleTestSuite encapsulates the tests
type SimpleTestSuite struct {
	suite.Suite
	ctx context.Context
	wpm *mocks.InterfaceManager
	wp  *mocks.InterfaceWorkerPool
	s   *SimpleExample
}

// SetupTest sets up often used objects
func (test *SimpleTestSuite) SetupTest() {
	test.T().Helper()

	test.ctx = context.TODO()
	test.wpm = mocks.NewInterfaceManager(test.T())
	test.wp = mocks.NewInterfaceWorkerPool(test.T())
	test.s = NewSimpleExample(test.wpm)
}

// TearDownTest tests whether the mock has been handled correctly after each test
func (test *SimpleTestSuite) TearDownTest() {
	test.wp.AssertExpectations(test.T())
	test.wpm.AssertExpectations(test.T())
}

// TestSimpleTestSuite Runs the testsuite
func TestSimpleTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SimpleTestSuite))
}

// TestNewSimpleExample tests the NewSimpleExample method
func (test *SimpleTestSuite) TestNewSimpleExample() {
	s := NewSimpleExample(test.wpm)
	test.NotNil(s)
	test.NotNil(s.wpm)
}

// TestSimplePool tests the SimplePool method
func (test *SimpleTestSuite) TestSimplePool() {
	wpm := workerpool.NewManager()
	s := NewSimpleExample(wpm)
	err := s.SimplePool(test.ctx, 50)
	test.NoError(err)
}

// TestSimplePool tests the SimplePool method
func (test *SimpleTestSuite) TestSimplePoolError() {
	wpm := workerpool.NewManager()
	s := NewSimpleExample(wpm)
	err := s.SimplePool(test.ctx, 100)
	test.NotNil(err)
	test.Contains(err.Error(), "Index cap reached!")
}

// TestSimplePoolMocked tests the SimplePool method
func (test *SimpleTestSuite) TestSimplePoolMocked() {
	test.wpm.EXPECT().NewWorkerPool(test.ctx).Return(test.wp).Once()
	test.wp.EXPECT().UseDefaults().Return(test.wp).Once()
	test.wp.EXPECT().Start().Return(test.wp).Once()
	test.wp.EXPECT().Quit().Return().Once()
	test.wp.EXPECT().AddTask(mock.AnythingOfType("*examples.SimpleTask")).Return(nil).Times(50)
	test.wp.EXPECT().Done().Return().Once()
	test.wp.EXPECT().Await().Return(nil).Once()

	err := test.s.SimplePool(test.ctx, 50)
	test.NoError(err)
}

// TestSimplePoolMocked tests the SimplePool method
func (test *SimpleTestSuite) TestSimplePoolMockedError() {
	expectedErr := fmt.Errorf("expected error")

	test.wpm.EXPECT().NewWorkerPool(test.ctx).Return(test.wp).Once()
	test.wp.EXPECT().UseDefaults().Return(test.wp).Once()
	test.wp.EXPECT().Start().Return(test.wp).Once()
	test.wp.EXPECT().Quit().Return().Once()
	test.wp.EXPECT().AddTask(mock.AnythingOfType("*examples.SimpleTask")).Return(nil).Times(74)
	test.wp.EXPECT().AddTask(mock.AnythingOfType("*examples.SimpleTask")).Return(expectedErr).Once()
	test.wp.EXPECT().Done().Return().Once()

	err := test.s.SimplePool(test.ctx, 100)
	test.NotNil(err)
	test.Equal(expectedErr.Error(), err.Error())
}
