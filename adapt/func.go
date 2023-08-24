package adapt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

func NewFuncAdaptor[
	T any,
	V Validatable[T],
	O any,
](
	domainCall func(context.Context, V) (O, error),
	decoder Decoder[T, V, O],
) (*FuncAdaptor[T, V, O], error) {
	if domainCall == nil {
		return nil, errors.New("cannot use a <nil> domain call")
	}
	var zero Decoder[T, V, O]
	if decoder == zero {
		return nil, errors.New("cannot use a <nil> decoder")
	}
	return &FuncAdaptor[T, V, O]{
		domainCall: domainCall,
		decoder:    decoder,
	}, nil
}

type FuncAdaptor[
	T any,
	V Validatable[T],
	O any,
] struct {
	domainCall func(context.Context, V) (O, error)
	decoder    Decoder[T, V, O]
}

func (a *FuncAdaptor[T, V, O]) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	request, encoder, err := a.decoder.Decode(w, r)
	if err != nil {
		return NewInvalidRequestError(fmt.Errorf("unable to decode: %w", err))
	}
	if err = request.Validate(); err != nil {
		return NewInvalidRequestError(err)
	}

	response, err := a.domainCall(r.Context(), request)
	if err != nil {
		return err
	}
	if err = encoder.Encode(w, response); err != nil {
		return fmt.Errorf("unable to encode: %w", err)
	}
	return nil
}

type StringFuncAdaptor[O any] struct {
	domainCall func(context.Context, string) (O, error)
	extractor  func(*http.Request) (string, error)
	encoder    Encoder[O]
}

func NewStringFuncAdaptor[O any](
	domainCall func(context.Context, string) (O, error),
	extractor func(*http.Request) (string, error),
	encoder Encoder[O],
) (*StringFuncAdaptor[O], error) {
	if domainCall == nil {
		return nil, errors.New("cannot use a <nil> domain call")
	}
	if extractor == nil {
		return nil, errors.New("cannot use a <nil> string extractor")
	}
	var zero Encoder[O]
	if encoder == zero {
		return nil, errors.New("cannot use a <nil> encoder")
	}
	return &StringFuncAdaptor[O]{
		domainCall: domainCall,
		extractor:  extractor,
		encoder:    encoder,
	}, nil
}

func (a *StringFuncAdaptor[O]) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	request, err := a.extractor(r)
	if err != nil {
		return NewInvalidRequestError(fmt.Errorf("unable to extract string: %w", err))
	}

	response, err := a.domainCall(r.Context(), request)
	if err != nil {
		return err
	}
	if err = a.encoder.Encode(w, response); err != nil {
		return fmt.Errorf("unable to encode: %w", err)
	}
	return nil
}
