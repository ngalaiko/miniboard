package httpx

import "net/http"

// Middleware used to define http middlewares.
type Middleware func(http.Handler) http.Handler
