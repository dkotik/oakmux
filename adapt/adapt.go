/*
Package adapt provides domain HTTP adaptors shaped to [oakmux.Handler].
*/
package adapt

import (
	"net/http"
)

// Validatable constrains a domain request. Validation errors will be wrapped as InvalidRequestError by the adapter.
type Validatable[T any] interface {
	*T
	Validate() error
}

type InvalidRequestError struct {
	error
}

func NewInvalidRequestError(fromError error) *InvalidRequestError {
	return &InvalidRequestError{fromError}
}

func (e *InvalidRequestError) Error() string {
	return "invalid request: " + e.error.Error()
}

func (e *InvalidRequestError) Unwrap() error {
	return e.error
}

func (e *InvalidRequestError) HyperTextStatusCode() int {
	return http.StatusUnprocessableEntity
}