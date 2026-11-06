package workerpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagerNewWorkerPool(t *testing.T) {
	t.Parallel()

	wp := NewManager().NewWorkerPool(t.Context())
	require.NotNil(t, wp)
	assert.Equal(t, defaultWorkerSize, wp.GetWorkers())
	assert.Equal(t, defaultBufferSize, wp.GetBufferSize())
}
