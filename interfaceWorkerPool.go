package workerpool

import (
	"github.com/LeoNdV001/workerpool/task"
)

// InterfaceWorkerPool is the interface implemented by WorkerPool
type InterfaceWorkerPool interface {
	// UseDefaults sets the default values for the WorkerPool
	//
	// Deprecated: NewWorkerPool already applies the defaults.
	UseDefaults() InterfaceWorkerPool
	// GetBufferSize returns the buffer size of the WorkerPool
	GetBufferSize() uint
	// WithBufferSize sets the buffer size of the WorkerPool
	WithBufferSize(bufferSize uint) InterfaceWorkerPool
	// GetWorkers returns the number of workers of the WorkerPool
	GetWorkers() uint
	// WithWorkers sets the workers size of the WorkerPool (minimum 1)
	WithWorkers(workers uint) InterfaceWorkerPool
	// Start launches the workers; a no-op while the pool is running
	Start() InterfaceWorkerPool
	// AddTask queues a task; errors if the pool is not running or its context is done
	AddTask(t task.InterfaceTask) error
	// Await blocks until all tasks have been processed (call Done first) or the
	// context is done, and returns the task errors joined
	Await() error
	// Done indicates that all tasks have been added
	Done()
	// Reset waits for queued tasks, clears collected errors and starts again.
	// Never call it from inside a task: it deadlocks.
	Reset()
	// Quit stops accepting tasks and blocks until all workers have exited.
	// Never call it from inside a task: it deadlocks.
	Quit()
}
