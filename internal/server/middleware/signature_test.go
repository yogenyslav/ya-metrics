package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/ya-metrics/pkg/secure"
)

func TestWithSignature(t *testing.T) {
	t.Parallel()

	key := "secure_key"
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})

	t.Run("valid signature", func(t *testing.T) {
		t.Parallel()

		body := []byte("test")

		sg := secure.NewSignatureGenerator(key)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		signature := sg.SignatureSHA256(body)
		req.Header.Set(headerSignature, signature)

		recorder := httptest.NewRecorder()
		signedHandler := WithSignature(key)(h)
		signedHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "test", recorder.Body.String())
	})

	t.Run("invalid signature", func(t *testing.T) {
		t.Parallel()

		body := []byte("test")

		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		req.Header.Set(headerSignature, "invalid_signature")

		recorder := httptest.NewRecorder()
		signedHandler := WithSignature(key)(h)
		signedHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("missing signature", func(t *testing.T) {
		t.Parallel()

		body := []byte("test")
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))

		recorder := httptest.NewRecorder()
		signedHandler := WithSignature(key)(h)
		signedHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("no key", func(t *testing.T) {
		t.Parallel()

		signedHandler := WithSignature("")(h)

		body := []byte("test")
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))

		recorder := httptest.NewRecorder()
		signedHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "test", recorder.Body.String())
	})
}
