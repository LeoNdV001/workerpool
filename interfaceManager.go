package workerpool

import (
	"context"
)

// InterfaceManager is the interface implemented by Manager
type InterfaceManager interface {
	// NewWorkerPool creates a new WorkerPool instance
	NewWorkerPool(ctx context.Context) InterfaceWorkerPool
}
