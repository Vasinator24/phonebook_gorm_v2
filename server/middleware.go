package server

import "net/http"

type Middleware func(http.Handler) http.Handler

type RouteGroup struct {
	mux         *http.ServeMux
	prefix      string
	middlewares []Middleware
}

func NewRouteGroup(
	mux *http.ServeMux,
	prefix string,
	middlewares ...Middleware,
) *RouteGroup {
	return &RouteGroup{
		mux:         mux,
		prefix:      prefix,
		middlewares: middlewares,
	}
}

func (g *RouteGroup) With(middlewares ...Middleware) *RouteGroup {
	groupMiddlewares := append([]Middleware{}, g.middlewares...)
	groupMiddlewares = append(groupMiddlewares, middlewares...)

	return NewRouteGroup(g.mux, g.prefix, groupMiddlewares...)
}

func (g *RouteGroup) Handle(path string, handler http.Handler) {
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		handler = g.middlewares[i](handler)
	}

	g.mux.Handle(g.prefix+path, handler)
}
