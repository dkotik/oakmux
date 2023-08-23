package adapt

import (
  "net/http"
  "encoding/json"
)

type Encoder[T any] interface {
  Encode(http.ResponseWriter, T) error
}

type Decoder[T any, V Validatable[T], O any] interface {
  Decode(http.ResponseWriter, *http.Request) (V, Encoder[O], error)
}

type Codec[T any, V Validatable[T], O any] interface {
  Encoder[O]
  Decoder[T, V, O]
}

type Finalizer[T any, V Validatable[T], O any] struct {
  Decoder Decoder[T, V, O]
  Finalize func(V) error
}

func (f *Finalizer[T, V, O]) Decode(
  w http.ResponseWriter,
  r *http.Request,
) (V, Encoder[O], error) {
  result, encoder, err := f.Decoder.Decode(w, r)
  if err != nil {
    return nil, nil, err
  }
  if err = f.Finalize(result); err != nil {
    return nil, nil, err
  }
  return result, encoder, nil
}

type jsonCodec[T any, V Validatable[T], O any] int64

func NewJSONCodec[T any, V Validatable[T], O any](maxBytes int64) Codec[T, V, O] {
  return jsonCodec[T, V, O](maxBytes)
}

func (j jsonCodec[T, V, O]) Encode(w http.ResponseWriter, value O) error {
  w.Header().Set("Content-Type", "application/json")
  return json.NewEncoder(w).Encode(&value)
}

func (j jsonCodec[T, V, O]) Decode(
  w http.ResponseWriter,
  r *http.Request,
) (V, Encoder[O], error) {
  var request V
  err := json.NewDecoder(http.MaxBytesReader(w, r.Body, int64(j))).Decode(&request)
  defer func() {
    r.Body.Close()
  }()
  if err != nil {
    return nil, nil, err
  }
  return request, j, nil
}
