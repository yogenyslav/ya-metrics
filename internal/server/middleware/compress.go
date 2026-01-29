package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/yogenyslav/ya-metrics/pkg/errs"
)

// GzipCompression is the compression type for gzip.
const GzipCompression string = "gzip"

var gzipWriterPool = &sync.Pool{
	New: func() any {
		return gzip.NewWriter(io.Discard)
	},
}

type writer interface {
	io.WriteCloser
	Reset(io.Writer)
}

type nopCloser struct {
	io.Writer
}

// Close implements io.Closer Close method.
func (c nopCloser) Close() error {
	return nil
}

// Reset implements gzip.Reset method.
func (c nopCloser) Reset(io.Writer) {}

// compressionResponseWriter is a wrapper around http.ResponseWriter to handle response compression.
type compressionResponseWriter struct {
	w               http.ResponseWriter
	compression     writer
	compressionType string
}

// newCompressResponseWriter creates a new compressionResponseWriter with specified compression type.
func newCompressResponseWriter(w http.ResponseWriter, compressionType string) *compressionResponseWriter {
	var compression writer

	switch compressionType {
	case GzipCompression:
		compression = gzipWriterPool.Get().(*gzip.Writer)
	default:
		compression = nopCloser{Writer: w}
	}

	return &compressionResponseWriter{
		w:               w,
		compression:     compression,
		compressionType: compressionType,
	}
}

// Header implements http.ResponseWriter Header method.
func (c *compressionResponseWriter) Header() http.Header {
	return c.w.Header()
}

// Write implements http.ResponseWriter Write method.
func (c *compressionResponseWriter) Write(b []byte) (int, error) {
	c.compression.Reset(c.w)
	return c.compression.Write(b)
}

// WriteHeader implements http.ResponseWriter WriteHeader method.
func (c *compressionResponseWriter) WriteHeader(statusCode int) {
	if statusCode >= 200 && statusCode < 300 {
		c.w.Header().Set("Content-Encoding", string(c.compressionType))
	}
	c.w.WriteHeader(statusCode)
}

// Close the compression writer.
func (c *compressionResponseWriter) Close() error {
	defer gzipWriterPool.Put(c.compression)
	return c.compression.Close()
}

// compressionReader is a wrapper around io.ReadCloser to handle request decompression.
type compressionReader struct {
	r               io.ReadCloser
	compression     io.ReadCloser
	compressionType string
}

// newCompressReader creates a new compressionReader with specified compression type.
func newCompressReader(r *http.Request, compressionType string) (*compressionReader, error) {
	var compression io.ReadCloser

	switch compressionType {
	case GzipCompression:
		gr, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		compression = gr
	default:
		compression = r.Body
	}

	return &compressionReader{
		r:               r.Body,
		compression:     compression,
		compressionType: compressionType,
	}, nil
}

// Read implements io.ReadCloser Read method.
func (c *compressionReader) Read(p []byte) (int, error) {
	return c.compression.Read(p)
}

// Close implements io.ReadCloser Close method.
func (c *compressionReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.compression.Close()
}

// WithCompression enables compression of specified type for HTTP requests.
func WithCompression(compression string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/debug/pprof") {
				next.ServeHTTP(w, r)
				return
			}

			acceptEncoding := r.Header.Get("Accept-Encoding")
			if strings.Contains(acceptEncoding, compression) {
				writer := newCompressResponseWriter(w, compression)
				defer writer.Close()
				w = writer
			}

			contentEncoding := r.Header.Get("Content-Encoding")
			if strings.Contains(contentEncoding, compression) {
				reader, err := newCompressReader(r, compression)
				if err != nil {
					http.Error(w, errs.Wrap(err, "create decompression reader").Error(), http.StatusInternalServerError)
					return
				}

				defer reader.Close()
				r.Body = reader
			}

			next.ServeHTTP(w, r)
		})
	}
}
