package core_http_middleware

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddlewares(
	handler http.Handler,
	middlewares ...Middleware,
) http.Handler {
	if len(middlewares) == 0 {
		return handler
	}

	for _, middleware := range slices.Backward(middlewares) {
		handler = middleware(handler)
	}

	return handler
}
