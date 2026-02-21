package repository

import "time"

// Ticker is an interface for time.Ticker to allow mocking in tests.
type Ticker interface {
	C() <-chan time.Time
	Stop()
}

// TickerFactory is a factory type for Ticker.
type TickerFactory func(d time.Duration) Ticker
