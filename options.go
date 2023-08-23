package oakmux

import (
	"errors"
	"fmt"
)

type options struct {
	redirectTrailingSlash bool // TODO: implement.
	handlers              map[*Route]Handler
	middleware            []Middleware
	prefix                string
	routes                map[string]*Route
	tree                  *Node
}

type Option func(*options) error

func WithRoute(name, pattern string, h Handler, mws ...Middleware) Option {
	return func(o *options) error {
		if name == "" {
			return fmt.Errorf("cannot use an empty route name")
		}
		if _, ok := o.routes[name]; ok {
			return fmt.Errorf("route %q is already set", name)
		}
		pattern = o.prefix + pattern
		if h == nil {
			return fmt.Errorf("cannot set an empty handler for path %q", pattern)
		}
		// p, err := Parse(pattern)
		// if err != nil {
		// 	return fmt.Errorf("cannot parse routing pattern %s: %w", pattern, err)
		// }
		// for kname, known := range o.routes {
		// 	relationship := comparePaths(p, known)
		// 	if relationship == equivalent || relationship == overlaps {
		// 		return fmt.Errorf("pattern %q conflicts with route pattern %s[%s]: %s",
		// 			p, kname, known, describeRel(p, known))
		// 	}
		// }
		// o.tree.addSegments(p.segments, p)
		// o.routes[name] = p

		route, err := NewRoute(pattern)
		if err != nil {
			return fmt.Errorf("cannot parse routing pattern %s: %w", pattern, err)
		}
		o.tree.Grow(route, route.segments)
		o.routes[name] = route
		if len(mws) > 0 {
			panic("middleware application disabled for now")
		}
		// o.handlers[p] = ApplyMiddleware(h, mws)
		o.handlers[route] = h
		return nil
	}
}

func WithMiddleware(mws ...Middleware) Option {
	return func(o *options) error {
		if len(mws) == 0 {
			return errors.New("WithMiddleware option requires at least one middleware")
		}
		for i, mw := range mws {
			if mw == nil {
				return fmt.Errorf("middleware %d is <nil>", len(o.middleware)+i)
			}
		}
		o.middleware = append(o.middleware, mws...)
		return nil
	}
}

func WithPrefix(p string) Option {
	return func(o *options) error {
		if p == "" {
			return errors.New("cannot use an empty route prefix")
		}
		if o.prefix != "" {
			return errors.New("route prefix is already set")
		}
		if len(o.routes) > 0 || len(o.handlers) > 0 {
			return errors.New("cannot set route prefix after routes have been added")
		}
		o.prefix = p
		return nil
	}
}

// func WithoutTrailingSlashRedirects() MuxOption {
// 	return func(o *muxOptions) error {
// 		if o.redirectTrailingSlash == false {
// 			return errors.New("trailing slash redirects are already disabled")
// 		}
// 		o.redirectTrailingSlash = true
// 		return nil
// 	}
// }
