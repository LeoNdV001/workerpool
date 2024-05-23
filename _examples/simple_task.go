package examples

import (
	"context"
	"fmt"

	"github.com/LeoNdV001/workerpool/task"
)

// SimpleTask is a simple task
type SimpleTask struct {
	Index uint
}

// NewSimpleTask creates a new SimpleTask instance
func NewSimpleTask(index uint) task.InterfaceTask {
	return &SimpleTask{
		Index: index,
	}
}

// Execute executes the task
func (st SimpleTask) Execute(_ context.Context) error {
	fmt.Printf("SimpleTask %d executing...\n", st.Index)

	if st.Index >= 75 {
		return fmt.Errorf("Index cap reached! (%d)\n", st.Index)
	}

	return nil
}

// OnSuccess is called when the task has been executed successfully
func (st SimpleTask) OnSuccess(_ context.Context) {
	fmt.Printf("SimpleTask %d executed successfully\n", st.Index)
}

// OnFailure is called when the task has failed
func (st SimpleTask) OnFailure(_ context.Context, err error) {
	fmt.Printf("SimpleTask %d failed with error %s\n", st.Index, err.Error())
}
