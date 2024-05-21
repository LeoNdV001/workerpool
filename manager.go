package workerpool

import (
	"context"
)

type Manager struct{}

// NewManager creates a new Manager instance
func NewManager() InterfaceManager {
	return &Manager{}
}

// NewWorkerPool creates a new WorkerPool instance
func (m *Manager) NewWorkerPool(ctx context.Context) InterfaceWorkerPool {
	return NewWorkerPool(ctx)
}
