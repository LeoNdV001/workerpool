package workerpool

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/LeoNdV001/workerpool/task"
)

const (
	defaultWorkerSize uint = 10
	defaultBufferSize uint = 100
)

var (
	// ErrResultsAlreadyRead is returned by AddTask when the pool does not accept
	// tasks: Done or Quit has been called, or Start has not been called yet.
	ErrResultsAlreadyRead = errors.New("workerpool: results have already been read")
	// ErrContextCanceled is returned when the pool's context ends before the work does.
	ErrContextCanceled = errors.New("workerpool: context canceled")
)

// WorkerPool is a collection of goroutines, where the number of concurrent
// goroutines processing requests does not exceed the specified maximum.
// Create it with NewWorkerPool; the zero value is not usable.
//
//nolint:containedctx // the pool is bound to one context for its lifetime
type WorkerPool struct {
	ctx             context.Context
	workers, buffer uint

	// mu guards tasks, closed and finished. AddTask holds the read lock while it
	// blocks on a send, so workers must never take mu: that is why errs has its
	// own lock. Merging the two locks deadlocks a full queue.
	mu       sync.RWMutex
	tasks    chan task.InterfaceTask
	closed   bool
	finished chan struct{} // closed when all workers of the current run have exited

	errMu sync.Mutex // guards errs
	// ponytail: keeps every task error until Await; count or cap them if a pool
	// sees millions of failures.
	errs []error
}

// NewWorkerPool creates a new WorkerPool instance with default settings
func NewWorkerPool(ctx context.Context) InterfaceWorkerPool {
	// A fresh pool behaves like one whose previous run has fully drained.
	finished := make(chan struct{})
	close(finished)

	return &WorkerPool{
		ctx:      ctx,
		workers:  defaultWorkerSize,
		buffer:   defaultBufferSize,
		closed:   true,
		finished: finished,
	}
}

// UseDefaults sets the default values for the WorkerPool
//
// Deprecated: NewWorkerPool already applies the defaults.
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

// WithWorkers sets the workers size of the WorkerPool (minimum 1)
func (w *WorkerPool) WithWorkers(workers uint) InterfaceWorkerPool {
	w.workers = max(1, workers)

	return w
}

// Start launches the workers. It is a no-op while the pool is running.
// After Done it waits for the previous run to drain, then starts a fresh one.
func (w *WorkerPool) Start() InterfaceWorkerPool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.closed {
		return w
	}

	<-w.finished

	w.errMu.Lock()
	w.errs = nil
	w.errMu.Unlock()

	var wg sync.WaitGroup

	tasks := make(chan task.InterfaceTask, w.buffer)

	for range w.workers {
		wg.Go(func() { w.work(tasks) })
	}

	finished := make(chan struct{})
	go func() {
		wg.Wait()
		close(finished)
	}()

	w.tasks, w.finished, w.closed = tasks, finished, false

	return w
}

// AddTask queues a task. It blocks while the buffer is full and returns an
// error if the pool does not accept tasks or its context is done.
func (w *WorkerPool) AddTask(t task.InterfaceTask) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.closed {
		return ErrResultsAlreadyRead
	}

	select {
	case w.tasks <- t:
		return nil
	case <-w.ctx.Done():
		return w.ctxErr()
	}
}

// Await blocks until all tasks have been processed or the context is done, and
// returns the task errors joined. Call Done first, or Await only returns once
// the context ends. If the context is done, the result also wraps
// ErrContextCanceled and the context's cause.
func (w *WorkerPool) Await() error {
	select {
	case <-w.running():
	case <-w.ctx.Done():
	}

	w.errMu.Lock()
	defer w.errMu.Unlock()

	if w.ctx.Err() != nil {
		return errors.Join(append([]error{w.ctxErr()}, w.errs...)...)
	}

	return errors.Join(w.errs...)
}

// Done indicates that all tasks have been added. Workers exit once the queue is empty.
func (w *WorkerPool) Done() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.closed {
		w.closed = true
		close(w.tasks)
	}
}

// Reset waits for the queued tasks to finish, clears the collected errors and
// starts again with the current settings. Call Await first to keep the errors.
// Never call it from inside a task: it waits for that task's worker and deadlocks.
func (w *WorkerPool) Reset() {
	w.Done()
	w.Start()
}

// Quit stops accepting tasks and blocks until all workers have exited.
// Never call it from inside a task: it waits for that task's worker and deadlocks.
func (w *WorkerPool) Quit() {
	w.Done()
	<-w.running()
}

// running returns the channel that is closed once the current run has finished.
func (w *WorkerPool) running() <-chan struct{} {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.finished
}

// work executes tasks until the queue is closed and empty, or the context is done.
func (w *WorkerPool) work(tasks <-chan task.InterfaceTask) {
	for {
		select {
		case <-w.ctx.Done():
			return
		case t, ok := <-tasks:
			if !ok {
				return
			}

			if err := t.Execute(w.ctx); err != nil {
				t.OnFailure(w.ctx, err)

				w.errMu.Lock()
				w.errs = append(w.errs, err)
				w.errMu.Unlock()

				continue
			}

			t.OnSuccess(w.ctx)
		}
	}
}

func (w *WorkerPool) ctxErr() error {
	return fmt.Errorf("%w: %w", ErrContextCanceled, context.Cause(w.ctx))
}
