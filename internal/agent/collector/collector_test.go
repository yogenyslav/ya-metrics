package collector

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestCollector_Collect(t *testing.T) {
	t.Parallel()

	t.Run(
		"Collect starts and stops correctly", func(t *testing.T) {
			t.Parallel()

			pollInterval := 1
			c := NewCollector(pollInterval, zerolog.Ctx(context.Background()))

			initialPollCount := c.GeneralMetrics().PollCount.Value
			initialAlloc := c.MemoryMetrics().Alloc.Value

			ctx, cancel := context.WithCancel(context.Background())
			c.Collect(ctx)

			<-time.After(time.Second * time.Duration(pollInterval+1))
			cancel()

			assert.Greater(t, c.GeneralMetrics().PollCount.Value, initialPollCount)
			assert.Greater(t, c.MemoryMetrics().Alloc.Value, initialAlloc)
		},
	)

	t.Run(
		"Collect stops on context cancel", func(t *testing.T) {
			t.Parallel()

			pollInterval := 1
			c := NewCollector(pollInterval, zerolog.DefaultContextLogger)
			initialPollCount := c.GeneralMetrics().PollCount.Value
			initialRandomValue := c.GeneralMetrics().RandomValue.Value
			initialAlloc := c.MemoryMetrics().Alloc.Value

			ctx, cancel := context.WithCancel(context.Background())
			c.Collect(ctx)

			<-time.After(time.Second * time.Duration(pollInterval/2))
			cancel()

			assert.Equal(t, initialPollCount, c.GeneralMetrics().PollCount.Value)
			assert.Equal(t, initialRandomValue, c.GeneralMetrics().RandomValue.Value)
			assert.Equal(t, initialAlloc, c.MemoryMetrics().Alloc.Value)
		},
	)
}
