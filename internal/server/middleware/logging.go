package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// LogData is a collective struct for logging HTTP requests.
type LogData struct {
	URI      string
	Method   string
	Duration time.Duration
}

// logResponseWriter is a wrapper around http.logResponseWriter to store status code and body size and log them.
type logResponseWriter struct {
	http.ResponseWriter

	StatusCode int
	BodySize   int
}

// WriteHeader stores the status code.
func (w *logResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write stores the size of the response body.
func (w *logResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.BodySize += size
	return size, err
}

// log the request data.
func (w *logResponseWriter) log(ctx context.Context, data *LogData) {
	zerolog.Ctx(ctx).
		Info().
		Str("method", data.Method).
		Str("uri", data.URI).
		Int("status", w.StatusCode).
		Int("size", w.BodySize).
		Dur("duration", data.Duration).
		Send()
}

// WithLogging enables logging middleware for HTTP requests.
func WithLogging(l *zerolog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &logResponseWriter{
				ResponseWriter: w,
			}

			start := time.Now()
			next.ServeHTTP(writer, r)
			duration := time.Since(start)

			data := &LogData{
				URI:      r.RequestURI,
				Method:   r.Method,
				Duration: duration,
			}
			writer.log(l.WithContext(r.Context()), data)
		})
	}
}
