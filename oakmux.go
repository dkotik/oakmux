/*
Package oakmux is a tree HTTP router with generic domain adaptors.
*/
package oakmux

import (
	"context"
	"fmt"
	"net/http"
)

type Error interface {
	error
	HyperTextStatusCode() int
}

type Handler interface {
	ServeHyperText(http.ResponseWriter, *http.Request) error
}

type Middleware func(Handler) Handler

// ApplyMiddleware applies [Middleware] to a [Handler] in reverse
// to preserve the logical order.
func ApplyMiddleware(h Handler, middleware ...Middleware) Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (f HandlerFunc) ServeHyperText(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type mux struct {
	handlers map[*Route]Handler
	routes   map[string]*Route
	tree     *Node
}

func New(withOptions ...Option) Handler {
	o := &options{
		handlers: make(map[*Route]Handler, 0),
		routes:   make(map[string]*Route),
		tree:     &Node{},
	}

	var err error
	for _, option := range withOptions {
		if err = option(o); err != nil {
			panic(fmt.Errorf("cannot initialize a path multiplexer: %w", err))
		}
	}

	return ApplyMiddleware(&mux{
		handlers: o.handlers,
		routes:   o.routes,
		tree:     o.tree,
	}, o.middleware...)
	return &mux{
		handlers: o.handlers,
		routes:   o.routes,
		tree:     o.tree,
	}
}

func (m *mux) ServeHyperText(w http.ResponseWriter, r *http.Request) error {
	route, matches := m.tree.MatchPath(r.URL.Path)
	// log.Println(route, matches, m.handlers)
	handler, ok := m.handlers[route]
	if !ok {
		return ErrNoRouteMatched
	}
	return handler.ServeHyperText(w, r.WithContext(
		context.WithValue(r.Context(), muxContextKey, &RoutingContext{
			mux:     m,
			matched: route,
			matches: matches,
		}),
	))
}

type noRouteMatchedError struct{}

func (e *noRouteMatchedError) Error() string {
	return "Not Found"
}

func (e *noRouteMatchedError) HyperTextStatusCode() int {
	return http.StatusNotFound
}

var ErrNoRouteMatched Error = &noRouteMatchedError{}
