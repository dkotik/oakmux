package oakmux

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dkotik/oakmux/adapt"
)

func WithPutHandler(h Handler, mws ...Middleware) MethodMuxOption {
	return func(o *methodMuxOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> get request handler")
		}
		if o.Put != nil {
			return errors.New("get request handler is already set")
		}
		o.Put = ApplyMiddleware(h, mws...)
		o.allowed += "," + http.MethodPut
		return nil
	}
}

func WithPutFunc[T any, V adapt.Validatable[T], O any](
	domainCall func(context.Context, V) (O, error),
	mws ...Middleware,
) MethodMuxOption {
	return func(o *methodMuxOptions) (err error) {
		adapted, err := adapt.NewUnaryFuncAdaptor(
			domainCall,
			adapt.NewJSONCodec[T, V, O](),
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for PUT method: %w", err)
		}
		return WithPutHandler(adapted, mws...)(o)
	}
}

func WithPutCustomFunc[T any, V adapt.Validatable[T], O any](
	domainCall func(context.Context, V) (O, error),
	codec adapt.Codec[T, V, O],
	mws ...Middleware,
) MethodMuxOption {
	return func(o *methodMuxOptions) (err error) {
		adapted, err := adapt.NewUnaryFuncAdaptor(
			domainCall,
			codec,
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for PUT method: %w", err)
		}
		return WithPutHandler(adapted, mws...)(o)
	}
}

func WithPutNullaryFunc[O any](
	domainCall func(context.Context) (O, error),
	mws ...Middleware,
) MethodMuxOption {
	return func(o *methodMuxOptions) (err error) {
		adapted, err := adapt.NewNullaryFuncAdaptor(
			domainCall,
			adapt.NewJSONEncoder[O](),
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for PUT method: %w", err)
		}
		return WithPutHandler(adapted, mws...)(o)
	}
}

func WithPutNullaryCustomFunc[O any](
	domainCall func(context.Context) (O, error),
	encoder adapt.Encoder[O],
	mws ...Middleware,
) MethodMuxOption {
	return func(o *methodMuxOptions) (err error) {
		adapted, err := adapt.NewNullaryFuncAdaptor(
			domainCall,
			encoder,
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for PUT method: %w", err)
		}
		return WithPutHandler(adapted, mws...)(o)
	}
}

func WithPutVoidFunc[T any, V adapt.Validatable[T]](
	domainCall func(context.Context, V) error,
	mws ...Middleware,
) MethodMuxOption {
	return func(o *methodMuxOptions) (err error) {
		adapted, err := adapt.NewVoidFuncAdaptor(
			domainCall,
			adapt.NewJSONCodec[T, V, T](),
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for PUT method: %w", err)
		}
		return WithPutHandler(adapted, mws...)(o)
	}
}

func WithPutCustomVoidFunc[T any, V adapt.Validatable[T]](
	domainCall func(context.Context, V) error,
	codec adapt.Codec[T, V, T],
	mws ...Middleware,
) MethodMuxOption {
	return func(o *methodMuxOptions) (err error) {
		adapted, err := adapt.NewVoidFuncAdaptor(
			domainCall,
			codec,
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for PUT method: %w", err)
		}
		return WithPutHandler(adapted, mws...)(o)
	}
}

func WithPutStringFunc[O any](
	domainCall func(context.Context, string) (O, error),
	extractor func(*http.Request) (string, error),
	mws ...Middleware,
) MethodMuxOption {
	return func(o *methodMuxOptions) error {
		adapted, err := adapt.NewStringUnaryFuncAdaptor(
			domainCall,
			extractor,
			adapt.NewJSONEncoder[O](),
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for PUT method: %w", err)
		}
		return WithPutHandler(adapted, mws...)(o)
	}
}

func WithPutCustomStringFunc[O any](
	domainCall func(context.Context, string) (O, error),
	extractor func(*http.Request) (string, error),
	encoder adapt.Encoder[O],
	mws ...Middleware,
) MethodMuxOption {
	return func(o *methodMuxOptions) error {
		adapted, err := adapt.NewStringUnaryFuncAdaptor(
			domainCall,
			extractor,
			encoder,
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for PUT method: %w", err)
		}
		return WithPutHandler(adapted, mws...)(o)
	}
}

func WithPutStringVoidFunc(
	domainCall func(context.Context, string) error,
	extractor func(*http.Request) (string, error),
	mws ...Middleware,
) MethodMuxOption {
	return func(o *methodMuxOptions) error {
		adapted, err := adapt.NewStringVoidFuncAdaptor(
			domainCall,
			extractor,
		)
		if err != nil {
			return fmt.Errorf("cannot adapt domain call for PUT method: %w", err)
		}
		return WithPutHandler(adapted, mws...)(o)
	}
}
