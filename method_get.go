package oakmux

import (
	"errors"
	"net/http"
)

func WithGet(h Handler) MethodMuxOption {
	return func(o *muxByMethodOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> get request handler")
		}
		if o.Get != nil {
			return errors.New("get request handler is already set")
		}
		o.Get = h
		o.allowed += "," + http.MethodGet
		return nil
	}
}

// TODO: add generic adaptors.
