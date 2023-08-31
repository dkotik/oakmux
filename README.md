# OakMux
## HTTP Router with Generic Domain Adaptors

## Inspiration

The core functionality was inspired by Jonathan Amsterdam's [alternative
ServeMux implementation][implementation] draft. For full discussion
see the [proposal here][proposal].

[implementation]: https://github.com/jba/muxpatterns/blob/main/tree.go
[proposal]: https://github.com/golang/go/discussions/60227

Mr. Amsterdam's draft contains unnecessary complexity for detecting route intersections, because it includes HTTP method differentiation inside the routing tree and does not anchor to the end of the request path without an explicit `[$]` segment.

In this implementation, all request paths are anchored to the end, unless ending with `[...]`, and the routing by method is set aside into a separate `MethodMux` handler. Therefore, the route intersections are always correctly caught at initialization, when the routing tree grows. This should also lead to slightly faster performance when working with live application routing trees.
