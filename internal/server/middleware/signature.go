package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/yogenyslav/ya-metrics/pkg/secure"
)

const headerSignature = "HashSHA256"

// SignatureGenerator is an interface for generating hash signatures.
type SignatureGenerator interface {
	SignatureSHA256(data []byte) string
}

// WithSignature is a middleware that checks incoming signatures of requests and adds signatures to outgoing responses.
func WithSignature(key string) Middleware {
	var sg SignatureGenerator
	if key != "" {
		sg = secure.NewSignatureGenerator(key)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			incomingSignature := r.Header.Get(headerSignature)
			if key == "" || incomingSignature == "" {
				next.ServeHTTP(w, r)
				return
			}

			body := bytes.Buffer{}
			_, err := io.Copy(&body, r.Body)
			if err != nil {
				http.Error(w, "failed to read request body", http.StatusInternalServerError)
				return
			}

			expectSignature := sg.SignatureSHA256(body.Bytes())
			if incomingSignature != expectSignature {
				http.Error(w, "invalid signature", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(&body)
			next.ServeHTTP(w, r)

			w.Header().Set(headerSignature, expectSignature)
		})
	}
}
