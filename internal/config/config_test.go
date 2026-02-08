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
		assert.Equal(t, defaultServerAddr, cfg.Server.Addr)
		assert.Equal(t, "debug", cfg.Server.LogLevel)
		assert.Equal(t, "metrics.json", cfg.Dump.FileStoragePath)
		assert.Equal(t, defaultStoreIntervalSec, cfg.Dump.StoreInterval)
	})

	t.Run("env vars are set", func(t *testing.T) {
		t.Setenv("ADDRESS", "localhost:9090")
		t.Setenv("LOG_LEVEL", "info")
		t.Setenv("FILE_STORAGE_PATH", "data.json")
		t.Setenv("STORE_INTERVAL", "600")

		cfg, err := NewConfig()
		require.NoError(t, err)
		assert.Equal(t, "localhost:9090", cfg.Server.Addr)
		assert.Equal(t, "info", cfg.Server.LogLevel)
		assert.Equal(t, "data.json", cfg.Dump.FileStoragePath)
		assert.Equal(t, 600, cfg.Dump.StoreInterval)
	})

	t.Run("Args are set", func(t *testing.T) {
		args := []string{
			"-a", "localhost:8080",
			"-l", "error",
			"-f", "metrics_data.json",
			"-i", "300",
		}

		os.Args = append([]string{"cmd"}, args...)
		cfg, err := NewConfig()
		require.NoError(t, err)
		assert.Equal(t, "localhost:8080", cfg.Server.Addr)
		assert.Equal(t, "error", cfg.Server.LogLevel)
		assert.Equal(t, "metrics_data.json", cfg.Dump.FileStoragePath)
		assert.Equal(t, 300, cfg.Dump.StoreInterval)
	})
}
