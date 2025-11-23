package retry_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ya-metrics/pkg/retry"
)

func TestWithLinearBackoffRetry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *retry.Config
		fn      func() error
		wantErr bool
	}{
		{
			name: "Success on first try",
			cfg: &retry.Config{
				MaxRetries:         3,
				LinearBackoffMilli: 100,
			},
			fn: func() error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "Success on third try",
			cfg: &retry.Config{
				MaxRetries:         3,
				LinearBackoffMilli: 100,
			},
			fn: func() func() error {
				cnt := 0
				return func() error {
					if cnt < 2 {
						cnt++
						return errors.New("err")
					}
					return nil
				}
			}(),
			wantErr: false,
		},
		{
			name: "Fail after max retries",
			cfg: &retry.Config{
				MaxRetries:         2,
				LinearBackoffMilli: 100,
			},
			fn: func() error {
				return errors.New("err")
			},
			wantErr: true,
		},
		{
			name: "Nil config",
			cfg:  nil,
			fn: func() error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "Nil config with error",
			cfg:  nil,
			fn: func() error {
				return errors.New("err")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := retry.WithLinearBackoffRetry(tt.cfg, tt.fn)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
