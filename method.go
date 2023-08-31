package oakmux

import (
	"errors"
	"fmt"
	"net/http"
)

type methodMux struct {
	Get     Handler
	Post    Handler
	Put     Handler
	Patch   Handler
	Delete  Handler
	allowed string
}

func (m *methodMux) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
		if m.Get != nil {
			return m.Get.ServeHyperText(w, r)
		}
	case http.MethodPost:
		if m.Post != nil {
			return m.Post.ServeHyperText(w, r)
		}
	case http.MethodPut:
		if m.Put != nil {
			return m.Put.ServeHyperText(w, r)
		}
	case http.MethodPatch:
		if m.Patch != nil {
			return m.Patch.ServeHyperText(w, r)
		}
	case http.MethodDelete:
		if m.Delete != nil {
			return m.Delete.ServeHyperText(w, r)
		}
	case http.MethodOptions:
		w.Header().Set("Allow", m.allowed)
		// w.WriteHeader(http.StatusOK)
		return nil
	}
	return NewMethodNotAllowedError(r.Method)
}

type methodNotAllowedError struct {
	method string
}

func NewMethodNotAllowedError(method string) Error {
	if method == "" {
		method = "unknown method"
	}
	return &methodNotAllowedError{method: method}
}

func (e *methodNotAllowedError) Error() string {
	return "method not allowed: " + e.method
}

func (e *methodNotAllowedError) HyperTextStatusCode() int {
	return http.StatusMethodNotAllowed
}

func NewMethodMux(withOptions ...MethodMuxOption) (Handler, error) {
	o := &methodMuxOptions{
		allowed: http.MethodOptions,
	}
	for _, option := range withOptions {
		if err := option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize method switch: %w", err)
		}
	}

	return &methodMux{
		Get:     o.Get,
		Post:    o.Post,
		Put:     o.Put,
		Patch:   o.Patch,
		Delete:  o.Delete,
		allowed: o.allowed,
	}, nil
}

type methodMuxOptions struct {
	Get     Handler
	Post    Handler
	Put     Handler
	Patch   Handler
	Delete  Handler
	allowed string
}

type MethodMuxOption func(*methodMuxOptions) error

func WithPost(h Handler) MethodMuxOption {
	return func(o *methodMuxOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> post request handler")
		}
		if o.Post != nil {
			return errors.New("post request handler is already set")
		}
		o.Post = h
		o.allowed += "," + http.MethodPost
		return nil
	}
}

func WithPut(h Handler) MethodMuxOption {
	return func(o *methodMuxOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> put request handler")
		}
		if o.Put != nil {
			return errors.New("put request handler is already set")
		}
		o.Put = h
		o.allowed += "," + http.MethodPut
		return nil
	}
}

func WithPatch(h Handler) MethodMuxOption {
	return func(o *methodMuxOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> patch request handler")
		}
		if o.Patch != nil {
			return errors.New("patch request handler is already set")
		}
		o.Patch = h
		o.allowed += "," + http.MethodPatch
		return nil
	}
}

func WithDelete(h Handler) MethodMuxOption {
	return func(o *methodMuxOptions) error {
		if h == nil {
			return errors.New("cannot use a <nil> delete request handler")
		}
		if o.Delete != nil {
			return errors.New("delete request handler is already set")
		}
		o.Delete = h
		o.allowed += "," + http.MethodDelete
		return nil
	}
}
