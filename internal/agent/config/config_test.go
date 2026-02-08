package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("Default values are used", func(t *testing.T) {
		cfg, err := NewConfig()
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080", cfg.ServerAddr)
		assert.Equal(t, defaultPollInterval, cfg.PollIntervalSec)
		assert.Equal(t, defaultReportInterval, cfg.ReportIntervalSec)
		assert.Equal(t, "", cfg.CompressionType)
		assert.Equal(t, "", cfg.SecureKey)
		assert.Equal(t, 1, cfg.RateLimit)
		assert.Equal(t, defaultBatchSize, cfg.BatchSize)
	})

	t.Run("env vars are set", func(t *testing.T) {
		t.Setenv("ADDRESS", "http://localhost:9090")
		t.Setenv("POLL_INTERVAL", "5")
		t.Setenv("REPORT_INTERVAL", "15")
		t.Setenv("COMPRESSION_TYPE", "gzip")
		t.Setenv("KEY", "securekey")
		t.Setenv("RATE_LIMIT", "10")

		cfg, err := NewConfig()
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:9090", cfg.ServerAddr)
		assert.Equal(t, 5, cfg.PollIntervalSec)
		assert.Equal(t, 15, cfg.ReportIntervalSec)
		assert.Equal(t, "gzip", cfg.CompressionType)
		assert.Equal(t, "securekey", cfg.SecureKey)
		assert.Equal(t, 10, cfg.RateLimit)
		assert.Equal(t, defaultBatchSize, cfg.BatchSize)
	})

	t.Run("Args are set", func(t *testing.T) {
		args := []string{
			"-a", "http://localhost:8080",
			"-p", "10",
			"-r", "20",
			"-c", "gzip",
			"-k", "securekey",
			"-l", "5",
		}

		os.Args = append([]string{"cmd"}, args...)
		cfg, err := NewConfig()
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080", cfg.ServerAddr)
		assert.Equal(t, 10, cfg.PollIntervalSec)
		assert.Equal(t, 20, cfg.ReportIntervalSec)
		assert.Equal(t, "gzip", cfg.CompressionType)
		assert.Equal(t, "securekey", cfg.SecureKey)
		assert.Equal(t, 5, cfg.RateLimit)
		assert.Equal(t, defaultBatchSize, cfg.BatchSize)
	})
}
