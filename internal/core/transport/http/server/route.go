package core_http_server

import (
	core_http_middleware "chatping-server/internal/core/transport/http/middleware"
	"net/http"
)

type Route struct {
	Method      string
	Path        string
	handler     http.HandlerFunc
	middlewares []core_http_middleware.Middleware
}

func NewRoute(
	method string,
	path string,
	handler http.HandlerFunc,
	middlewares ...core_http_middleware.Middleware,
) Route {
	return Route{
		Method:      method,
		Path:        path,
		handler:     handler,
		middlewares: middlewares,
	}
}

func (r *Route) Handler() http.Handler {
	return core_http_middleware.ChainMiddlewares(
		r.handler,
		r.middlewares...,
	)
}
