package httpx

import (
	"net/http"
)

// Middleware used to define http middlewares.
type Middleware func(http.Handler) http.Handler

// Chain applies all middlewares to the handler. All middlewares will be called in the order they are defined.
func Chain(handler http.Handler, mm ...Middleware) http.Handler {
	wrappedHandler := handler
	for i := len(mm) - 1; i >= 0; i-- {
		wrappedHandler = mm[i](wrappedHandler)
	}
	return wrappedHandler
}
