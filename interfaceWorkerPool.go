package workerpool

import (
	"github.com/LeoNdV001/workerpool/task"
)

// InterfaceWorkerPool is the interface implemented by WorkerPool
type InterfaceWorkerPool interface {
	// UseDefaults sets the default values for the WorkerPool
	UseDefaults() InterfaceWorkerPool
	// GetBufferSize returns the buffer size of the WorkerPool
	GetBufferSize() uint
	// WithBufferSize sets the buffer size of the WorkerPool
	WithBufferSize(bufferSize uint) InterfaceWorkerPool
	// GetWorkers returns the number of workers of the WorkerPool
	GetWorkers() uint
	// WithWorkers sets the workers size of the WorkerPool
	WithWorkers(workers uint) InterfaceWorkerPool
	// Start initializes the WorkerPool for processing tasks
	Start() InterfaceWorkerPool
	// AddTask adds a tasks to the WorkerPool
	AddTask(t task.InterfaceTask) error
	// Await waits until all tasks have been processed. Make sure to
	// call the done function in order for the error channel to be closed.
	Await() error
	// Done indicate that all tasks have been published and close the tasks channel.
	Done()
	// Reset resets the WorkerPool to its initial state
	Reset()
	// Quit gracefully shuts down the WorkerPool
	Quit()
}
