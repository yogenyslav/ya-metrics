package middleware

import "net/http"

// Middleware is the type for HTTP middleware functions.
type Middleware func(next http.Handler) http.Handler
