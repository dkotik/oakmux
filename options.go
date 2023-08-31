package oakmux

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dkotik/oakmux/adapt"
)

const DefaultRequestReadLimitOf1MB = 1 << 20

type options struct {
	redirectTrailingSlash bool // TODO: implement.
	limitlessRequestBytes bool
	maximumRequestBytes   int64
	handlers              map[*Route]Handler
	middleware            []Middleware
	prefix                string
	routes                map[string]*Route
	tree                  *Node
}

type Option func(*options) error

func WithRequestReadLimitOf(maximumBytes int64) Option {
	return func(o *options) error {
		if maximumBytes <= 0 {
			return errors.New("read limit must be greater than 0 bytes")
		}
		if o.maximumRequestBytes != 0 {
			return fmt.Errorf("read limit is already set to: %d", o.maximumRequestBytes)
		}
		if o.limitlessRequestBytes {
			return errors.New("cannot override WithLimitlessRequest option")
		}
		o.maximumRequestBytes = maximumBytes
		return nil
	}
}

func WithDefaultRequestReadLimitOf1MB() Option {
	return func(o *options) error {
		if o.maximumRequestBytes != 0 || o.limitlessRequestBytes {
			return nil // already set
		}
		if err := WithRequestReadLimitOf(DefaultRequestReadLimitOf1MB)(o); err != nil {
			return fmt.Errorf("unable to set default maximum JSON request bytes: %w", err)
		}
		return nil
	}
}

func WithLimitlessRequestBytes() Option {
	return func(o *options) error {
		if o.maximumRequestBytes != 0 {
			return fmt.Errorf("read limit is already set to: %d", o.maximumRequestBytes)
		}
		if o.limitlessRequestBytes {
			return errors.New("read limit is already removed")
		}
		o.limitlessRequestBytes = true
		return nil
	}
}

func WithRouteHandler(name, pattern string, h Handler, mws ...Middleware) Option {
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

		route, err := NewRoute(name, pattern)
		if err != nil {
			return fmt.Errorf("cannot parse routing pattern %s: %w", pattern, err)
		}
		if err = o.tree.Grow(route, route.segments); err != nil {
			return fmt.Errorf("cannot use routing pattern %s for route %s: %w", pattern, name, err)
		}
		o.routes[name] = route
		o.handlers[route] = ApplyMiddleware(h, mws...)
		return nil
	}
}

func WithRouteFunc[T any, V adapt.Validatable[T], O any](
	name, pattern string,
	domainCall func(context.Context, V) (O, error),
	mws ...Middleware,
) Option {
	return func(o *options) (err error) {
		adapted, err := adapt.NewUnaryFuncAdaptor(
			domainCall,
			adapt.NewJSONCodec[T, V, O](),
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for route %q at path %q: %w", name, pattern, err)
		}
		return WithRouteHandler(name, pattern, adapted, mws...)(o)
	}
}

func WithRouteCustomFunc[T any, V adapt.Validatable[T], O any](
	name, pattern string,
	domainCall func(context.Context, V) (O, error),
	codec adapt.Codec[T, V, O],
	mws ...Middleware,
) Option {
	return func(o *options) (err error) {
		adapted, err := adapt.NewUnaryFuncAdaptor(
			domainCall,
			codec,
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for route %q at path %q: %w", name, pattern, err)
		}
		return WithRouteHandler(name, pattern, adapted, mws...)(o)
	}
}

func WithRouteNullaryFunc[O any](
	name, pattern string,
	domainCall func(context.Context) (O, error),
	mws ...Middleware,
) Option {
	return func(o *options) (err error) {
		adapted, err := adapt.NewNullaryFuncAdaptor(
			domainCall,
			adapt.NewJSONEncoder[O](),
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for route %q at path %q: %w", name, pattern, err)
		}
		return WithRouteHandler(name, pattern, adapted, mws...)(o)
	}
}

func WithRouteNullaryCustomFunc[O any](
	name, pattern string,
	domainCall func(context.Context) (O, error),
	encoder adapt.Encoder[O],
	mws ...Middleware,
) Option {
	return func(o *options) (err error) {
		adapted, err := adapt.NewNullaryFuncAdaptor(
			domainCall,
			encoder,
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for route %q at path %q: %w", name, pattern, err)
		}
		return WithRouteHandler(name, pattern, adapted, mws...)(o)
	}
}

func WithRouteVoidFunc[T any, V adapt.Validatable[T]](
	name, pattern string,
	domainCall func(context.Context, V) error,
	mws ...Middleware,
) Option {
	return func(o *options) (err error) {
		adapted, err := adapt.NewVoidFuncAdaptor(
			domainCall,
			adapt.NewJSONCodec[T, V, T](),
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for route %q at path %q: %w", name, pattern, err)
		}
		return WithRouteHandler(name, pattern, adapted, mws...)(o)
	}
}

func WithRouteCustomVoidFunc[T any, V adapt.Validatable[T]](
	name, pattern string,
	domainCall func(context.Context, V) error,
	codec adapt.Codec[T, V, T],
	mws ...Middleware,
) Option {
	return func(o *options) (err error) {
		adapted, err := adapt.NewVoidFuncAdaptor(
			domainCall,
			codec,
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for route %q at path %q: %w", name, pattern, err)
		}
		return WithRouteHandler(name, pattern, adapted, mws...)(o)
	}
}

func WithRouteStringFunc[O any](
	name, pattern string,
	domainCall func(context.Context, string) (O, error),
	extractor func(*http.Request) (string, error),
	mws ...Middleware,
) Option {
	return func(o *options) error {
		adapted, err := adapt.NewStringUnaryFuncAdaptor(
			domainCall,
			extractor,
			adapt.NewJSONEncoder[O](),
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for route %q at path %q: %w", name, pattern, err)
		}
		return WithRouteHandler(name, pattern, adapted, mws...)(o)
	}
}

func WithRouteCustomStringFunc[O any](
	name, pattern string,
	domainCall func(context.Context, string) (O, error),
	extractor func(*http.Request) (string, error),
	encoder adapt.Encoder[O],
	mws ...Middleware,
) Option {
	return func(o *options) error {
		adapted, err := adapt.NewStringUnaryFuncAdaptor(
			domainCall,
			extractor,
			encoder,
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for route %q at path %q: %w", name, pattern, err)
		}
		return WithRouteHandler(name, pattern, adapted, mws...)(o)
	}
}

func WithRouteStringVoidFunc(
	name, pattern string,
	domainCall func(context.Context, string) error,
	extractor func(*http.Request) (string, error),
	mws ...Middleware,
) Option {
	return func(o *options) error {
		adapted, err := adapt.NewStringVoidFuncAdaptor(
			domainCall,
			extractor,
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for route %q at path %q: %w", name, pattern, err)
		}
		return WithRouteHandler(name, pattern, adapted, mws...)(o)
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
