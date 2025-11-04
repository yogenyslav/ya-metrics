package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCompressionHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Helper()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			require.NoError(t, err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}

func TestWithCompression(t *testing.T) {
	t.Parallel()

	data := "test data"
	gzipData := func() []byte {
		buf := &bytes.Buffer{}
		w := gzip.NewWriter(buf)
		w.Write([]byte(data))
		w.Close()
		return buf.Bytes()
	}

	tests := []struct {
		name        string
		compression string
		body        []byte
		headers     http.Header
		wantBody    string
		wantHeaders http.Header
	}{
		{
			name:     "No compression",
			body:     []byte(data),
			wantBody: data,
		},
		{
			name:        "Gzip request compression",
			compression: GzipCompression,
			body:        gzipData(),
			headers: http.Header{
				"Content-Encoding": []string{GzipCompression},
			},
			wantBody: data,
		},
		{
			name:        "Gzip response compression",
			compression: GzipCompression,
			body:        []byte(data),
			headers: http.Header{
				"Accept-Encoding": []string{GzipCompression},
			},
			wantBody: string(gzipData()),
			wantHeaders: http.Header{
				"Content-Encoding": []string{GzipCompression},
			},
		},
		{
			name:        "Gzip request and response compression",
			compression: GzipCompression,
			body:        gzipData(),
			headers: http.Header{
				"Content-Encoding": []string{GzipCompression},
				"Accept-Encoding":  []string{GzipCompression},
			},
			wantBody: string(gzipData()),
			wantHeaders: http.Header{
				"Content-Encoding": []string{GzipCompression},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := WithCompression(tt.compression)(testCompressionHandler(t))
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(tt.body))

			if tt.headers != nil {
				req.Header = tt.headers
			}

			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, req)

			assert.Equal(t, tt.wantBody, recorder.Body.String())
			for k, wantValues := range tt.wantHeaders {
				v := recorder.Header().Values(k)
				assert.ElementsMatch(t, wantValues, v)
			}
		})
	}
}
