package examples

import (
	"context"

	"github.com/LeoNdV001/workerpool"
)

func example() error {
	wp := workerpool.NewWorkerPool(context.Background()).UseDefaults().Start()
	defer wp.Quit()

	for i := range uint(10) {
		_ = wp.AddTask(NewSimpleTask(i))
	}

	wp.Done()

	return wp.Await()
}
