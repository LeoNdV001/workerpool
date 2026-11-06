package workerpool

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errBoom = errors.New("boom")

// recorder collects task outcomes from concurrent workers
type recorder struct {
	mu        sync.Mutex
	succeeded []string
	failed    int
}

// testTask sleeps for delay, then returns err
type testTask struct {
	name  string
	err   error
	delay time.Duration
	rec   *recorder
}

func (t *testTask) Execute(_ context.Context) error {
	time.Sleep(t.delay)

	return t.err
}

func (t *testTask) OnSuccess(_ context.Context) {
	t.rec.mu.Lock()
	defer t.rec.mu.Unlock()

	t.rec.succeeded = append(t.rec.succeeded, t.name)
}

func (t *testTask) OnFailure(_ context.Context, _ error) {
	t.rec.mu.Lock()
	defer t.rec.mu.Unlock()

	t.rec.failed++
}

func TestProcessesAllTasks(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pool            func(InterfaceWorkerPool) InterfaceWorkerPool
		workers, buffer uint
	}{
		"defaults":   {func(p InterfaceWorkerPool) InterfaceWorkerPool { return p }, 10, 100},
		"custom":     {func(p InterfaceWorkerPool) InterfaceWorkerPool { return p.WithBufferSize(10).WithWorkers(5) }, 5, 10},
		"unbuffered": {func(p InterfaceWorkerPool) InterfaceWorkerPool { return p.WithBufferSize(0).WithWorkers(0) }, 1, 0},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			rec := &recorder{}
			wp := tt.pool(NewWorkerPool(t.Context())).Start()
			defer wp.Quit()

			assert.Equal(t, tt.workers, wp.GetWorkers())
			assert.Equal(t, tt.buffer, wp.GetBufferSize())

			go func() {
				for i := range 500 {
					assert.NoError(t, wp.AddTask(&testTask{name: fmt.Sprint(i), rec: rec}))
				}

				wp.Done()
			}()

			require.NoError(t, wp.Await())
			assert.Len(t, rec.succeeded, 500)
		})
	}
}

// Failing tasks must not block workers when nobody is in Await yet.
func TestCollectsErrorsWithoutConcurrentAwait(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		rec := &recorder{}
		wp := NewWorkerPool(t.Context()).WithWorkers(2).WithBufferSize(0).Start()
		defer wp.Quit()

		for i := range 100 {
			var err error
			if i%10 == 0 {
				err = errBoom
			}

			require.NoError(t, wp.AddTask(&testTask{name: fmt.Sprint(i), err: err, rec: rec}))
		}

		wp.Done()

		err := wp.Await()
		require.ErrorIs(t, err, errBoom)
		assert.Equal(t, 10, strings.Count(err.Error(), "boom"))
		assert.Equal(t, 10, rec.failed)
		assert.Len(t, rec.succeeded, 90)
	})
}

func TestContextCanceled(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()

		wp := NewWorkerPool(ctx).WithBufferSize(0).Start()
		defer wp.Quit()

		// Done is never called, so only the deadline can end Await.
		err := wp.Await()
		require.ErrorIs(t, err, ErrContextCanceled)
		require.ErrorIs(t, err, context.DeadlineExceeded)

		// Workers are gone; AddTask must fail instead of blocking forever.
		require.ErrorIs(t, wp.AddTask(&testTask{rec: &recorder{}}), ErrContextCanceled)
	})
}

func TestAddTaskWhenNotRunning(t *testing.T) {
	t.Parallel()

	wp := NewWorkerPool(t.Context())
	require.ErrorIs(t, wp.AddTask(&testTask{}), ErrResultsAlreadyRead, "before Start")

	wp.Start().Done()
	require.ErrorIs(t, wp.AddTask(&testTask{}), ErrResultsAlreadyRead, "after Done")
}

func TestResetDrainsAndRestarts(t *testing.T) {
	t.Parallel()

	rec := &recorder{}
	wp := NewWorkerPool(t.Context()).WithBufferSize(0).WithWorkers(0).Start()
	defer wp.Quit()

	for i := range 500 {
		require.NoError(t, wp.AddTask(&testTask{name: fmt.Sprint(i), rec: rec}))
	}

	wp.Reset()
	require.NoError(t, wp.AddTask(&testTask{name: "after-reset", rec: rec}))
	wp.Done()

	require.NoError(t, wp.Await())
	require.Len(t, rec.succeeded, 501)
	assert.Equal(t, "0", rec.succeeded[0])
	assert.Equal(t, "after-reset", rec.succeeded[500])
}

func TestResetClearsErrors(t *testing.T) {
	t.Parallel()

	wp := NewWorkerPool(t.Context()).Start()
	defer wp.Quit()

	require.NoError(t, wp.AddTask(&testTask{err: errBoom, rec: &recorder{}}))
	wp.Reset()
	wp.Done()

	require.NoError(t, wp.Await())
}

func TestQuitWaitsForWorkers(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		rec := &recorder{}
		wp := NewWorkerPool(t.Context()).WithWorkers(3).Start()

		for range 9 {
			require.NoError(t, wp.AddTask(&testTask{delay: time.Second, rec: rec}))
		}

		start := time.Now()
		wp.Quit()

		assert.Equal(t, 3*time.Second, time.Since(start))
		assert.Len(t, rec.succeeded, 9)

		wp.Quit() // second call is a no-op
	})
}

// Run with -race: AddTask racing Done must never panic with "send on closed channel".
func TestConcurrentAddTaskAndDone(t *testing.T) {
	t.Parallel()

	rec := &recorder{}
	wp := NewWorkerPool(t.Context()).WithBufferSize(1).Start()

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 100 {
				err := wp.AddTask(&testTask{rec: rec})
				if err != nil {
					assert.ErrorIs(t, err, ErrResultsAlreadyRead)

					return
				}
			}
		})
	}

	wp.Done()
	wg.Wait()
	require.NoError(t, wp.Await())
}
