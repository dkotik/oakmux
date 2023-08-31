# OakMux: HTTP Router with Generic Adaptors

## Inspiration

The core functionality was inspired by Jonathan Amsterdam's [alternative
ServeMux implementation][implementation] draft. For full discussion
see the [proposal here][proposal].

[implementation]: https://github.com/jba/muxpatterns/blob/main/tree.go
[proposal]: https://github.com/golang/go/discussions/60227

Mr. Amsterdam's draft contains unnecessary complexity for detecting route intersections, because it includes HTTP method differentiation inside the routing tree and does not anchor to the end of the request path without an explicit `[$]` segment.

In this implementation, all request paths are anchored to the end, unless ending with `[...]`, and the routing by method is set aside into a separate `MethodMux` handler. Therefore, the route intersections are always correctly caught at initialization, when the routing tree grows. This should also lead to slightly faster performance when working with live application routing trees.

## Domain Adaptors

Domain logic adaptors come in three general flavors:

1. UnaryFunc: func(context, inputStruct) (outputStruct, error)
2. NullaryFunc: func(context) (outputStruct, error)
3. VoidFunc: func(context, inputStruct) error

Each input requires implementation of `adapt.Validatable` for safety. Validation errors are decorated with the correct `http.StatusUnprocessableEntity` status code.
