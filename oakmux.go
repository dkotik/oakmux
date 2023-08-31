/*
Package oakmux is a tree HTTP router with generic domain adaptors.
*/
package oakmux

import (
	"context"
	"fmt"
	"net/http"
)

func Must[T any](this T, err error) T {
	if err != nil {
		panic(err)
	}
	return this
}

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
		// TODO: check for nils?
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

func New(withOptions ...Option) (Handler, error) {
	o := &options{
		redirectToTrailingSlash:   true,
		redirectFromTrailingSlash: true,
		handlers:                  make(map[*Route]Handler, 0),
		routes:                    make(map[string]*Route),
		tree:                      &Node{},
	}

	var err error
	for _, option := range append(
		withOptions,
		WithDefaultRequestReadLimitOf1MB(),
		func(o *options) error {
			if o.limitlessRequestBytes {
				return nil
			}
			o.middleware = append([]Middleware{ // inject read limiting middleware
				NewRequestReadLimiterMiddleware(o.maximumRequestBytes),
			}, o.middleware...)
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize a path multiplexer: %w", err)
		}
	}

	if err = o.injectTrailingSlashRedirects(); err != nil {
		return nil, fmt.Errorf("cannot add a trailing slash redirect: %w", err)
	}

	if len(o.middleware) > 0 {
		return ApplyMiddleware(&mux{
			handlers: o.handlers,
			routes:   o.routes,
			tree:     o.tree,
		}, o.middleware...), nil
	}
	return &mux{
		handlers: o.handlers,
		routes:   o.routes,
		tree:     o.tree,
	}, nil
}

func (m *mux) ServeHyperText(w http.ResponseWriter, r *http.Request) error {
	route, matches := m.tree.MatchPath(r.URL.Path)
	handler, ok := m.handlers[route]
	if !ok {
		return ErrNoRouteMatched
	}
	// log.Println(route, matches, m.handlers)
	// log.Println(route.String(), r.URL.Path)
	return handler.ServeHyperText(w, r.WithContext(
		context.WithValue(r.Context(), muxContextKey, &RoutingContext{
			mux:     m,
			matched: route,
			matches: matches,
		}),
	))
}

func (m *mux) String() string {
	return m.tree.String()
}
