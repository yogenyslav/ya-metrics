package retry

import (
	"errors"
	"time"
)

// ErrUnretriable is an error that can't be retried.
var ErrUnretriable = errors.New("unretriable error")

var (
	// DefaultRetries is the default number of retries.
	DefaultRetries int = 3
	// DefaultLinearBackoffMilli is the default linear backoff in milliseconds.
	DefaultLinearBackoffMilli int = 2000
)

// Config holds retry configuration.
type Config struct {
	MaxRetries         int
	LinearBackoffMilli int
}

// WithLinearBackoffRetry is a wrapper for retry logic with linear backoff.
func WithLinearBackoffRetry(cfg *Config, fn func() error) error {
	var err error

	if cfg == nil {
		return fn()
	}

	for i := 0; i <= cfg.MaxRetries; i++ {
		err = fn()

		if err == nil {
			return nil
		}

		if errors.Is(err, ErrUnretriable) {
			return err
		}

		time.Sleep(time.Millisecond * time.Duration(cfg.LinearBackoffMilli*i))
	}

	return err
}
