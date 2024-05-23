package workerpool

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/LeoNdV001/workerpool/task"
)

// WorkerPoolTestSuite encapsulates the tests
type WorkerPoolTestSuite struct {
	suite.Suite
	ctx   context.Context
	tasks []task.InterfaceTask
}

// TestTask is a simple task object
type TestTask struct {
	Name    string
	Err     error
	Results *[]string
	T       *testing.T
}

// Execute handles the business logic of the task
func (t *TestTask) Execute(_ context.Context) error {
	return t.Err
}

// OnSuccess is called when the task has been executed successfully
func (t *TestTask) OnSuccess(_ context.Context) {
	if t.Results != nil {
		*t.Results = append(*t.Results, t.Name)
	}
}

// OnFailure is called when the task has failed
func (t *TestTask) OnFailure(_ context.Context, _ error) {}

// NewTask creates a new task object
func (test *WorkerPoolTestSuite) NewTask(name string, err error) task.InterfaceTask {
	return &TestTask{
		Name: name,
		Err:  err,
		T:    test.T(),
	}
}

// SetupTest sets up often used objects (for every test)
func (test *WorkerPoolTestSuite) SetupTest() {
	test.T().Helper()
	test.ctx = context.TODO()
	test.tasks = make([]task.InterfaceTask, 0)

	for i := range 500 {
		test.tasks = append(test.tasks, test.NewTask(fmt.Sprintf("task-%d\n", i), nil))
	}
}

// TearDownTest tests whether the mock has been handled correctly after each test
func (test *WorkerPoolTestSuite) TearDownTest() {
	os.Clearenv()
	test.ctx.Done()
}

// TestWorkerPoolTestSuite Runs the testsuite
func TestWorkerPoolTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WorkerPoolTestSuite))
}

// TestNewDefaultWorkerPool tests the NewWorkerPool function with a simple task object
func (test *WorkerPoolTestSuite) TestNewDefaultWorkerPool() {
	ctx, c := context.WithTimeout(test.ctx, 5*time.Second)
	defer c()

	workerPool := NewWorkerPool(ctx).UseDefaults().Start()
	test.NotNil(workerPool)

	defer workerPool.Quit()

	go func() {
		for i := range test.tasks {
			err := workerPool.AddTask(test.tasks[i])
			test.Nil(err)
		}

		// added all the jobs
		workerPool.Done()
	}()

	test.Nil(workerPool.Await())
}

// TestWorkerPoolOptions tests the WorkerPool with custom options
func (test *WorkerPoolTestSuite) TestWorkerPoolOptions() {
	ctx, c := context.WithTimeout(test.ctx, 5*time.Second)
	defer c()

	workerPool := NewWorkerPool(ctx).
		WithBufferSize(10).
		WithWorkers(5).
		Start()

	defer workerPool.Quit()

	test.NotNil(workerPool)
	test.Equal(uint(10), workerPool.GetBufferSize())
	test.Equal(uint(5), workerPool.GetWorkers())

	go func() {
		for i := range test.tasks {
			err := workerPool.AddTask(test.tasks[i])
			test.Nil(err)
		}

		// added all the jobs
		workerPool.Done()
	}()

	test.Nil(workerPool.Await())
}

// TestWorkerPoolUnbufferedSingleWorker tests the WorkerPool with custom options
func (test *WorkerPoolTestSuite) TestWorkerPoolUnbufferedSingleWorker() {
	ctx, c := context.WithTimeout(test.ctx, 5*time.Second)
	defer c()

	workerPool := NewWorkerPool(ctx).
		WithBufferSize(0).
		WithWorkers(0).
		Start()

	defer workerPool.Quit()

	test.NotNil(workerPool)
	test.Equal(uint(0), workerPool.GetBufferSize())
	test.Equal(uint(1), workerPool.GetWorkers())

	go func() {
		for i := range test.tasks {
			err := workerPool.AddTask(test.tasks[i])
			test.Nil(err)
		}

		// added all the jobs
		workerPool.Done()
	}()

	test.Nil(workerPool.Await())
}

// TestWorkerPoolWorkerReset tests the WorkerPool with custom options
func (test *WorkerPoolTestSuite) TestWorkerPoolWorkerResetAndResults() {
	ctx, c := context.WithTimeout(test.ctx, 5*time.Second)
	defer c()

	workerPool := NewWorkerPool(ctx).
		WithBufferSize(0).
		WithWorkers(0).
		Start()

	defer workerPool.Quit()

	test.NotNil(workerPool)
	test.Equal(uint(0), workerPool.GetBufferSize())
	test.Equal(uint(1), workerPool.GetWorkers())

	results := make([]string, 0)

	for i := range test.tasks {
		t := test.tasks[i].(*TestTask)
		t.Results = &results

		err := workerPool.AddTask(t)
		test.Nil(err)
	}

	workerPool.Reset()

	err := workerPool.AddTask(test.tasks[0])
	test.Nil(err)

	workerPool.Done()

	test.Nil(workerPool.Await())
	test.Len(results, 501)
	test.Equal("task-0\n", results[0])
	test.Equal("task-0\n", results[len(results)-1])
}

// TestWorkerPoolOptions tests the WorkerPool with custom options
func (test *WorkerPoolTestSuite) TestWorkerPoolOnFailure() {
	var (
		expectedError = fmt.Errorf("expected error")
		err, taskErr  error
	)

	ctx, c := context.WithTimeout(test.ctx, 5*time.Second)
	defer c()

	workerPool := NewWorkerPool(ctx).UseDefaults().Start()
	test.NotNil(workerPool)

	defer workerPool.Quit()

	go func() {
		for i := range 500 {
			err = nil
			if i%50 == 0 {
				err = expectedError
			}

			taskErr = workerPool.AddTask(test.NewTask(fmt.Sprintf("task-%d\n", i), err))
			test.Nil(taskErr)
		}

		// added all the jobs
		workerPool.Done()
	}()

	err = workerPool.Await()
	test.NotNil(err)
	test.Equal("expected error\nexpected error\nexpected error\nexpected error\nexpected error\n"+
		"expected error\nexpected error\nexpected error\nexpected error\nexpected error", err.Error())
}

// TestWorkerPoolContextCanceled tests the WorkerPool with a canceled context
func (test *WorkerPoolTestSuite) TestWorkerPoolContextCanceled() {
	ctx, c := context.WithTimeout(test.ctx, time.Millisecond)
	defer c()

	workerPool := NewWorkerPool(ctx).UseDefaults().Start()
	test.NotNil(workerPool)

	go func() {
		err := workerPool.AddTask(test.NewTask(fmt.Sprintf("task-%d\n", 1), nil))
		test.Nil(err)
	}()

	err := workerPool.Await()
	test.NotNil(err)
	test.Contains(err.Error(), ErrContextCanceled)
}

// TestWorkerPoolCantAddJobToFinishedPool tests the WorkerPool with a canceled context
func (test *WorkerPoolTestSuite) TestWorkerPoolCantAddJobToFinishedPool() {
	ctx, c := context.WithTimeout(test.ctx, 5*time.Second)
	defer c()

	workerPool := NewWorkerPool(ctx).UseDefaults().Start()
	test.NotNil(workerPool)

	workerPool.Done()

	err := workerPool.AddTask(test.NewTask(fmt.Sprintf("task-%d\n", 1), nil))
	test.NotNil(err)
	test.Contains(err.Error(), ErrResultsAlreadyRead)
}
