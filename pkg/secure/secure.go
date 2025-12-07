package secure

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"sync"
)

// SignatureGenerator is a struct for generating hash signatures.
type SignatureGenerator struct {
	key        []byte
	hashSha256 hash.Hash
	mu         *sync.Mutex
}

// NewSignatureGenerator creates a new SignatureGenerator with the given key and hash function.
// If no hash function is provided, it defaults to HMAC-SHA256.
func NewSignatureGenerator(key string) *SignatureGenerator {
	return &SignatureGenerator{
		key:        []byte(key),
		hashSha256: hmac.New(sha256.New, []byte(key)),
		mu:         &sync.Mutex{},
	}
}

// SignatureSHA256 generates a SHA256 HMAC signature for the given data.
func (sg *SignatureGenerator) SignatureSHA256(data []byte) string {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	sg.hashSha256.Reset()
	sg.hashSha256.Write(data)
	return hex.EncodeToString(sg.hashSha256.Sum(nil))
}
