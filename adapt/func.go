package adapt

import (
	"context"
	"fmt"
	"net/http"
)

type FuncAdaptor[
	T any,
	V Validatable[T],
	O any,
] struct {
	DomainCall func(context.Context, V) (O, error)
	Decoder    Decoder[T, V, O]
}

func (a *FuncAdaptor[T, V, O]) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	request, encoder, err := a.Decoder.Decode(w, r)
	if err != nil {
		return NewInvalidRequestError(fmt.Errorf("unable to decode: %w", err))
	}
	if err = request.Validate(); err != nil {
		return NewInvalidRequestError(err)
	}

	response, err := a.DomainCall(r.Context(), request)
	if err != nil {
		return err
	}
	if err = encoder.Encode(w, response); err != nil {
		return fmt.Errorf("unable to encode: %w", err)
	}
	return nil
}
