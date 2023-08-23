package oakmux

import (
	"context"
	"fmt"
	"strconv"
)

type muxContextKeyType struct{}

var muxContextKey = &muxContextKeyType{}

func GetRoutingContext(ctx context.Context) *RoutingContext {
	routing, _ := ctx.Value(muxContextKey).(*RoutingContext)
	return routing
}

type RoutingContext struct {
	mux     *mux
	matched *Route
	matches []string
}

func (r *RoutingContext) Path(routeName string, fields map[string]string) (string, error) {
	route, ok := r.mux.routes[routeName]
	if !ok {
		return "", ErrPathNotFound
	}
	return route.Path(fields)
}

func (r *RoutingContext) MatchedFields() *MatchedFields {
	bindings := make(map[string]string)
	i := 0
	// TODO: should loop over namedSegments instead of all segments.
	// for _, segment := range r.matched.segments {
	// 	if segment.wild || segment.multi {
	// 		bindings[segment.s] = r.matches[i]
	// 		i++
	// 	}
	// }
	for _, segment := range r.matched.namedSegments {
		bindings[segment.Name()] = r.matches[i]
		i++
	}
	return &MatchedFields{
		route:    r.matched,
		bindings: bindings,
	}
}

type MatchedFields struct {
	route    *Route
	bindings map[string]string
}

func (m *MatchedFields) Str(name string, value *string) error {
	fieldValue, ok := m.bindings[name]
	if !ok {
		return fmt.Errorf("route pattern %q does not contain field named %q", m.route, name)
	}
	*value = fieldValue
	return nil
}

func (m *MatchedFields) Int(name string, value *int) error {
	var s string
	if err := m.Str(name, &s); err != nil {
		return err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return ErrNotInteger
	}
	*value = i
	return nil
}

func (m *MatchedFields) Int64(name string, value *int64) error {
	var i int
	if err := m.Int(name, &i); err != nil {
		return err
	}
	*value = int64(i)
	return nil
}

func (m *MatchedFields) Uint(name string, value *uint) error {
	var i int
	if err := m.Int(name, &i); err != nil {
		return err
	}
	if i < 0 {
		return ErrNotUnsignedInteger
	}
	*value = uint(i)
	return nil
}

func (m *MatchedFields) Uint64(name string, value *uint64) error {
	var i int
	if err := m.Int(name, &i); err != nil {
		return err
	}
	if i < 0 {
		return ErrNotUnsignedInteger
	}
	*value = uint64(i)
	return nil
}

func (m *MatchedFields) Page(name string, value *int) error {
	var i int
	if err := m.Int(name, &i); err != nil {
		return err
	}
	if i < 1 {
		return ErrNotPageNumber
	}
	*value = i
	return nil
}
