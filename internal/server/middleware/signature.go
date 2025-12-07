package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/yogenyslav/ya-metrics/pkg/secure"
)

const headerSignature = "Hash"

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
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			body := bytes.Buffer{}
			_, err := io.Copy(&body, r.Body)
			if err != nil {
				http.Error(w, "failed to read request body", http.StatusInternalServerError)
				return
			}

			incomingSignature := r.Header.Get(headerSignature)
			expectSignature := sg.SignatureSHA256(body.Bytes())
			if incomingSignature != expectSignature {
				http.Error(w, "invalid signature", http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)

			w.Header().Set(headerSignature, expectSignature)
		})
	}
}
