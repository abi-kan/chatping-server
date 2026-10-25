package core_http_server

import (
	core_http_middleware "chatping-server/internal/core/transport/http/middleware"
	"fmt"
	"net/http"
)

type APIVersion string

var (
	APIVersion1 = APIVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion  APIVersion
	middlewares []core_http_middleware.Middleware
}

func NewAPIVersionRouter(
	apiVersion APIVersion,
	middlewares ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:    http.NewServeMux(),
		apiVersion:  apiVersion,
		middlewares: middlewares,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, route.Handler())
	}
}

func (r *APIVersionRouter) Handler() http.Handler {
	return core_http_middleware.ChainMiddlewares(
		r,
		r.middlewares...,
	)
}
