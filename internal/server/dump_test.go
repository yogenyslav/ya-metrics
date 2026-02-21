package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/yogenyslav/ya-metrics/internal/config"
	"github.com/yogenyslav/ya-metrics/internal/server/repository"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
)

type testTicker struct {
	ch chan time.Time
}

func (t *testTicker) C() <-chan time.Time {
	return t.ch
}

func (t *testTicker) Stop() {}

type mockDumper struct {
	mock.Mock

	called chan struct{}
}

func (d *mockDumper) Dump(ctx context.Context, gaugeRepo repository.Repo, counterRepo repository.Repo) error {
	args := d.Called(ctx, gaugeRepo, counterRepo)
	d.called <- struct{}{}
	return args.Error(0)
}

func TestServer_Dumping(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := &testTicker{ch: make(chan time.Time)}
	getTicker := func(time.Duration) Ticker { return ticker }
	called := make(chan struct{}, 1)
	dumper := &mockDumper{called: called}
	counterRepo := new(mocks.MockMetricRepo[int64])
	gaugeRepo := new(mocks.MockMetricRepo[float64])

	dumper.On("Dump", mock.Anything, gaugeRepo, counterRepo).Return(nil)

	s := &Server{
		cfg: &config.Config{
			Dump: &config.DumpConfig{
				StoreInterval: 1,
			},
		},
		dumper: dumper,
	}
	s.Dumping(ctx, dumper, getTicker, gaugeRepo, counterRepo)

	ticker.ch <- time.Now()

	select {
	case <-called:
		// ok
	case <-time.After(time.Second):
		t.Error("dump was not called within expected interval")
	}
}
