package examples

import (
	"context"

	"github.com/LeoNdV001/workerpool"
)

// SimpleExample is the method receiver
type SimpleExample struct {
	wpm workerpool.InterfaceManager
}

// NewSimpleExample initializes the service
func NewSimpleExample(wpm workerpool.InterfaceManager) *SimpleExample {
	return &SimpleExample{
		wpm: wpm,
	}
}

// SimplePool creates a new WorkerPool instance
func (e *SimpleExample) SimplePool(ctx context.Context, numberOfTasks uint) error {
	wp := e.wpm.NewWorkerPool(ctx).UseDefaults().Start()
	defer wp.Quit()

	for i := range numberOfTasks {
		t := NewSimpleTask(i)

		if err := wp.AddTask(t); err != nil {
			// Close channel
			wp.Done()

			return err
		}
	}

	// All tasks have been added
	wp.Done()

	return wp.Await()
}
