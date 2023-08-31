/*
Package adapt provides domain logic adaptors shaped to [oakmux.Handler]. Adaptors come in three flavors:

1. UnaryFunc: func(context, inputStruct) (outputStruct, error)
2. NullaryFunc: func(context) (outputStruct, error)
3. VoidFunc: func(context, inputStruct) error

Each input requires implementation of [Validatable] for safety. Validation errors are decorated with the correct [http.StatusUnprocessableEntity] status code.
*/
package adapt

import (
	"net/http"
)

// Validatable constrains a domain request. Validation errors are wrapped as [InvalidRequestError] by the adapter.
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
