package collector

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollector_Collect(t *testing.T) {
	t.Parallel()

	t.Run(
		"Collect starts and stops correctly", func(t *testing.T) {
			t.Parallel()

			pollInterval := 1
			c := NewCollector(pollInterval)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			c.Collect(ctx)
			initialPollCount := c.PollCount.Value
			initialRandomValue := c.RandomValue.Value
			initialAlloc := c.MemoryMetrics.Alloc.Value

			<-time.After(time.Second*time.Duration(pollInterval) + 500*time.Millisecond)

			assert.Greater(t, c.PollCount.Value, initialPollCount)
			assert.NotEqual(t, c.RandomValue.Value, initialRandomValue)
			assert.Greater(t, c.MemoryMetrics.Alloc.Value, initialAlloc)
		},
	)

	t.Run(
		"Collect stops on context cancel", func(t *testing.T) {
			t.Parallel()

			pollInterval := 1
			c := NewCollector(pollInterval)
			initialPollCount := c.PollCount.Value
			initialRandomValue := c.RandomValue.Value
			initialAlloc := c.MemoryMetrics.Alloc.Value

			ctx, cancel := context.WithCancel(context.Background())
			c.Collect(ctx)

			<-time.After(time.Second * time.Duration(pollInterval/2))
			cancel()

			assert.Equal(t, initialPollCount, c.PollCount.Value)
			assert.Equal(t, initialRandomValue, c.RandomValue.Value)
			assert.Equal(t, initialAlloc, c.MemoryMetrics.Alloc.Value)
		},
	)
}
