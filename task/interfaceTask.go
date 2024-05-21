package task

import "context"

// InterfaceTask is the interface for a task that can be executed by a worker
type InterfaceTask interface {
	// Execute performs the work
	Execute(ctx context.Context) error
	// OnSuccess handles the result of Execute()
	OnSuccess(ctx context.Context)
	// OnFailure handles any error returned from Execute()
	OnFailure(ctx context.Context, err error)
}
