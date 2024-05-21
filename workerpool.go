package workerpool

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"github.com/LeoNdV001/workerpool/task"
)

const (
	defaultWorkerSize     uint = 10
	defaultBufferSize     uint = 100
	ErrResultsAlreadyRead      = "WorkerPool: results have already been read"
	ErrContextCanceled         = "WorkerPool: context canceled"
)

// WorkerPool is a collection of goroutines, where the number of concurrent
// goroutines processing requests does not exceed the specified maximum.
//
//nolint:containedctx // context is set knowingly
type WorkerPool struct {
	ctx         context.Context
	tasks       chan task.InterfaceTask
	errors      chan error
	wg          *sync.WaitGroup
	hasFinished atomic.Bool
	workers,
	buffer uint
}

// NewWorkerPool creates a new WorkerPool instance
func NewWorkerPool(ctx context.Context) InterfaceWorkerPool {
	return &WorkerPool{
		ctx: ctx,
		wg:  &sync.WaitGroup{},
	}
}

// UseDefaults sets the default values for the WorkerPool
func (w *WorkerPool) UseDefaults() InterfaceWorkerPool {
	w.workers = defaultWorkerSize
	w.buffer = defaultBufferSize

	return w
}

// GetBufferSize returns the buffer size of the WorkerPool
func (w *WorkerPool) GetBufferSize() uint {
	return w.buffer
}

// WithBufferSize sets the buffer size of the WorkerPool
func (w *WorkerPool) WithBufferSize(bufferSize uint) InterfaceWorkerPool {
	w.buffer = bufferSize

	return w
}

// GetWorkers returns the number of workers of the WorkerPool
func (w *WorkerPool) GetWorkers() uint {
	return w.workers
}

// WithWorkers sets the workers size of the WorkerPool
func (w *WorkerPool) WithWorkers(workers uint) InterfaceWorkerPool {
	w.workers = uint(math.Max(1, float64(workers)))

	return w
}

// Start initializes the WorkerPool for processing tasks
func (w *WorkerPool) Start() InterfaceWorkerPool {
	w.resetChannels()
	w.resetWorkers()

	return w
}

// AddTask adds a tasks to the WorkerPool
func (w *WorkerPool) AddTask(t task.InterfaceTask) error {
	if w.hasFinished.Load() {
		return fmt.Errorf(ErrResultsAlreadyRead)
	}

	w.tasks <- t

	return nil
}

// Await waits until all tasks have been processed. Make sure to
// call the done function in order for the error channel to be closed.
func (w *WorkerPool) Await() error {
	errorList := make([]error, 0)

	for {
		select {
		case v, ok := <-w.errors:
			if !ok {
				w.errors = nil
			} else {
				errorList = append(errorList, v)
			}

		case <-w.ctx.Done():
			return fmt.Errorf("%s; %w", ErrContextCanceled, errors.Join(errorList...))
		}

		if w.errors == nil {
			break
		}
	}

	if len(errorList) > 0 {
		return errors.Join(errorList...)
	}

	return nil
}

// Done indicate that all tasks have been published and close the tasks channel.
func (w *WorkerPool) Done() {
	if !w.hasFinished.Load() {
		w.hasFinished.Store(true)

		close(w.tasks)
	}
}

// Reset resets the WorkerPool to its initial state
func (w *WorkerPool) Reset() {
	w.Done()
	w.resetChannels()
	w.resetWorkers()
}

// Quit gracefully shuts down the WorkerPool
func (w *WorkerPool) Quit() {
	w.Done()
	go w.awaitProcessingDone()
}

// resetChannels resets the channels
func (w *WorkerPool) resetChannels() {
	w.tasks = make(chan task.InterfaceTask, w.buffer)
	w.errors = make(chan error)
}

// resetWorkers resets the workers
func (w *WorkerPool) resetWorkers() {
	w.spawnWorkers()
	w.setHasFinished(false)

	go w.awaitProcessingDone()
}

// spawnWorkers spawns the workers
func (w *WorkerPool) spawnWorkers() {
	for range w.workers {
		w.wg.Add(1)
		go w.worker()
	}
}

// worker handles the execution of tasks
func (w *WorkerPool) worker() {
	defer w.wg.Done()

	for {
		select {
		case t, ok := <-w.tasks:
			if !ok {
				return
			}

			if err := t.Execute(w.ctx); err != nil {
				t.OnFailure(w.ctx, err)

				w.errors <- err

				continue
			}

			t.OnSuccess(w.ctx)
		case <-w.ctx.Done():
			return
		}

		if w.errors == nil {
			break
		}
	}
}

// awaitProcessingDone waits until all tasks have been processed
func (w *WorkerPool) awaitProcessingDone() {
	w.wg.Wait()
	w.setHasFinished(true)

	// close the errors channel if not already closed
	if w.errors != nil {
		close(w.errors)
	}
}

// setHasFinished sets the hasFinished flag
func (w *WorkerPool) setHasFinished(hasFinished bool) {
	w.hasFinished.Store(hasFinished)
}
