package oakmux

import (
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrDoubleSlash        = NewRoutingError(errors.New("path contains double slash"))
	ErrPathNotFound       = NewRoutingError(errors.New("path not found"))
	ErrNotInteger         = NewRoutingError(errors.New("field is not an integer"))
	ErrNotUnsignedInteger = NewRoutingError(errors.New("field is not an unsigned integer"))
	ErrNotPageNumber      = NewRoutingError(errors.New("field is not an page number"))
)

type RoutingError struct {
	cause error
}

func NewRoutingError(cause error) *RoutingError {
	return &RoutingError{cause: cause}
}

func (e *RoutingError) Unwrap() error {
	return e.cause
}

func (e *RoutingError) HyperTextStatusCode() int {
	return http.StatusNotFound
}

func (e *RoutingError) Error() string {
	return http.StatusText(http.StatusNotFound)
}

func (e *RoutingError) LogValue() slog.Value {
	return slog.StringValue("routing error: " + e.cause.Error())
}
