# workerpool

A small Go worker pool: a fixed number of goroutines run tasks from a buffered
queue, and `Await` returns every task error joined together.

Requires Go 1.27.2 or newer.

```bash
go get github.com/LeoNdV001/workerpool
```

## Usage

A task implements `task.InterfaceTask`:

```go
type InterfaceTask interface {
	Execute(ctx context.Context) error        // does the work
	OnSuccess(ctx context.Context)            // called when Execute returns nil
	OnFailure(ctx context.Context, err error) // called when Execute returns an error
}
```

Run tasks:

```go
wp := workerpool.NewWorkerPool(ctx). // defaults: 10 workers, buffer of 100
	WithWorkers(5).
	WithBufferSize(10).
	Start()
defer wp.Quit() // waits until all workers have exited

for _, t := range tasks {
	if err := wp.AddTask(t); err != nil { // blocks while the buffer is full
		return err
	}
}

wp.Done()         // no more tasks
err := wp.Await() // waits for the queue to drain
```

- `Await` returns `errors.Join` of all task errors. If `ctx` ends first, the
  error also matches `errors.Is(err, workerpool.ErrContextCanceled)` and the
  context's cause (for example `context.DeadlineExceeded`).
- `AddTask` returns `ErrResultsAlreadyRead` after `Done`/`Quit`, or before `Start`.
- `Reset` waits for queued tasks, clears collected errors and starts again.
- Never call `Quit` or `Reset` from inside a task: both wait for every worker,
  including the one running that task, so they deadlock.

## Testing with mocks

`NewManager()` returns an `InterfaceManager` whose `NewWorkerPool` you can inject
and replace in tests. Mockery mocks for `InterfaceManager`, `InterfaceWorkerPool`
and `task.InterfaceTask` live in `github.com/LeoNdV001/workerpool/_mocks`. See
[`_examples/`](_examples) for both a real and a mocked pool.

## Development

```bash
make test   # go test -race, root module and _examples
make lint   # golangci-lint v2
mockery     # regenerate _mocks from .mockery.yaml
```
