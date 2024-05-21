package workerpool

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ManagerTestSuite encapsulates the tests
type ManagerTestSuite struct {
	suite.Suite
	ctx context.Context
}

// SetupTest sets up often used objects (for every test)
func (test *ManagerTestSuite) SetupTest() {
	test.T().Helper()
	test.ctx = context.TODO()
}

// TearDownTest tests whether the mock has been handled correctly after each test
func (test *ManagerTestSuite) TearDownTest() {
	os.Clearenv()
	test.ctx.Done()
}

// TestWorkerPoolTestSuite Runs the testsuite
func TestManagerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ManagerTestSuite))
}

// TestNewManager tests the creation of a new Manager
func (test *ManagerTestSuite) TestNewManager() {
	m := NewManager()
	test.NotNil(m)
}

// TestNewWorkerPool tests the creation of a new WorkerPool
func (test *ManagerTestSuite) TestNewWorkerPool() {
	m := NewManager()
	test.NotNil(m)

	wp := m.NewWorkerPool(test.ctx)
	test.NotNil(wp)

	wp.UseDefaults()
	test.Equal(uint(10), wp.GetWorkers())
}
