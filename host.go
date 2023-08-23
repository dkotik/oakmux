package oakmux

import (
	"errors"
	"fmt"
	"net/http"

	"log/slog"
)

type hostMuxOptions struct {
	hosts    []string
	handlers []Handler
}

type HostMuxOption func(*hostMuxOptions) error

func WithHostHandler(host string, handler Handler) HostMuxOption {
	return func(o *hostMuxOptions) error {
		if host == "" {
			return errors.New("cannot use an empty host name")
		}
		if handler == nil {
			return errors.New("cannot use a <nil> handler")
		}
		for _, known := range o.hosts {
			if known == host {
				return fmt.Errorf("host %q already has a handler", known)
			}
		}
		o.hosts = append(o.hosts, host)
		o.handlers = append(o.handlers, handler)
		return nil
	}
}

// NewHostMux creates a [Handler] that multiplexes by [http.Request] host name.
func NewHostMux(withOptions ...HostMuxOption) (Handler, error) {
	var (
		o   = &hostMuxOptions{}
		err error
	)
	for _, option := range append(
		withOptions,
		func(o *hostMuxOptions) error {
			if len(o.hosts) == 0 || len(o.handlers) == 0 {
				return errors.New("empty host handler list")
			}
			return nil
		},
	) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot create a new host multiplexer: %w", err)
		}
	}

	// mapHostMux will be faster than list at 8 entries
	if len(o.hosts) >= 8 {
		mux := make(mapHostMux)
		for i, name := range o.hosts {
			mux[name] = o.handlers[i]
		}
		return mux, nil
	}
	return &listHostMux{
		hosts:    o.hosts,
		handlers: o.handlers,
	}, nil
}

type UnknownHostError struct {
	host string
}

func (e *UnknownHostError) Error() string {
	return http.StatusText(http.StatusNotFound)
}

func (e *UnknownHostError) HyperTextStatusCode() int {
	return http.StatusNotFound
}

func (e *UnknownHostError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("message", "uknown host"),
		slog.String("host", e.host),
	)
}

type mapHostMux map[string]Handler

func (h mapHostMux) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	name := r.URL.Hostname()
	handler, ok := h[name]
	if !ok {
		return &UnknownHostError{host: name}
	}
	return handler.ServeHyperText(w, r)
}

type listHostMux struct {
	hosts    []string
	handlers []Handler
}

func (l *listHostMux) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	name := r.URL.Hostname()
	for i, host := range l.hosts {
		if host == name {
			return l.handlers[i].ServeHyperText(w, r)
		}
	}
	return &UnknownHostError{host: name}
}
