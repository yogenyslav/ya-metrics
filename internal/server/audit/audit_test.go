package audit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/ya-metrics/internal/config"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("No sources", func(t *testing.T) {
		t.Parallel()

		cfg := &config.AuditConfig{}
		audit := New(cfg)
		assert.Len(t, audit.sources, 0)
	})

	t.Run("File source", func(t *testing.T) {
		t.Parallel()

		cfg := &config.AuditConfig{
			File: "audit.log",
		}
		audit := New(cfg)
		assert.Len(t, audit.sources, 1)
	})

	t.Run("File and service sources", func(t *testing.T) {
		t.Parallel()

		cfg := &config.AuditConfig{
			File: "audit.log",
			URL:  "http://audit-service.local/logs",
		}
		audit := New(cfg)
		assert.Len(t, audit.sources, 2)
	})
}

func Test_fileSource_Log(t *testing.T) {
	t.Parallel()

	t.Run("Log to file", func(t *testing.T) {
		t.Parallel()

		src := &fileSource{
			filePath: t.TempDir() + "/test_audit.log",
			mu:       &sync.Mutex{},
		}

		data := []byte(`{"ts":1625077765,"metrics":["metric1","metric2"],"ip_address":"127.0.0.1"}`)
		err := src.Log(context.Background(), data)
		assert.NoError(t, err)
	})

	t.Run("Log to invalid file path", func(t *testing.T) {
		t.Parallel()

		src := &fileSource{
			filePath: "",
			mu:       &sync.Mutex{},
		}

		data := []byte(`{"ts":1625077765,"metrics":["metric1","metric2"],"ip_address":"127.0.0.1"}`)
		err := src.Log(context.Background(), data)
		assert.Error(t, err)
	})
}

func Test_serviceSource_Log(t *testing.T) {
	t.Parallel()

	t.Run("Log to empty URL", func(t *testing.T) {
		t.Parallel()

		src := &serviceSource{
			url:    "",
			client: http.DefaultClient,
		}

		data := []byte(`{"ts":1625077765,"metrics":["metric1","metric2"],"ip_address":"127.0.0.1"}`)
		err := src.Log(context.Background(), data)
		assert.Error(t, err)
	})

	t.Run("Log to valid URL", func(t *testing.T) {
		t.Parallel()

		recorder := httptest.NewRecorder()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}))
		defer server.Close()

		src := &serviceSource{
			url:    server.URL,
			client: http.DefaultClient,
		}

		data := []byte(`{"ts":1625077765,"metrics":["metric1","metric2"],"ip_address":"127.0.0.1"}`)
		err := src.Log(context.Background(), data)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}

func TestAudit_LogMetrics(t *testing.T) {
	t.Parallel()

	t.Run("Log metrics to all sources, no errors", func(t *testing.T) {
		t.Parallel()

		recorder := httptest.NewRecorder()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}))
		defer server.Close()

		cfg := &config.AuditConfig{
			File: t.TempDir() + "/test_audit.log",
			URL:  server.URL,
		}
		audit := New(cfg)

		err := audit.LogMetrics(context.Background(), []string{"metric1", "metric2"}, "127.0.0.1")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("Log metrics with file source error", func(t *testing.T) {
		t.Parallel()

		cfg := &config.AuditConfig{
			File: "///",
		}
		audit := New(cfg)

		err := audit.LogMetrics(context.Background(), []string{"metric1", "metric2"}, "127.0.0.1")
		assert.Error(t, err)
	})
}
