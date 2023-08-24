package oakmux

import (
	"fmt"
	"strings"
)

type Route struct {
	name          string
	segments      []Segment
	namedSegments []Segment
}

func munchPath(p string) (
	segmentDefinition string,
	remainder string,
	err error,
) {
	if len(p) <= 1 {
		return p, "", nil
	}
	if p[0] == '/' {
		p = p[1:] // strip leading slash
	}
	// TODO: use IndexByte in all segment parsing.
	switch i := strings.IndexByte(p, '/'); i {
	case -1: // slash not found
		return p, "", nil
	case 0: // double slash
		return "", "", ErrDoubleSlash
	default:
		return p[:i], p[i:], nil
	}
}

func NewRoute(name, fromPath string) (r *Route, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("routing path %q is invalid: %w", fromPath, err)
		}
	}()

	r = &Route{
		name: name,
	}
	remainder := fromPath
	segmentDefinition := ""
	var currentType, lastSegmentType SegmentType

	for {
		segmentDefinition, remainder, err = munchPath(remainder)
		if err != nil {
			return nil, err
		}
		if segmentDefinition == "" {
			break
		}
		// fmt.Println("? segment: ", segmentDefinition, remainder)
		s, err := NewSegment(segmentDefinition)
		if err != nil {
			return nil, err
		}

		currentType = s.Type()
		if (lastSegmentType == SegmentTypeTerminal) && (currentType != SegmentTypeTrailingSlash) {
			return nil, fmt.Errorf("a terminal segment is followed by a %s segment", currentType)
		}
		lastSegmentType = currentType

		if currentType == SegmentTypeDynamic || (currentType == SegmentTypeTerminal && s.Name() != "") {
			name := s.Name()
			for _, s := range r.namedSegments {
				if name != "" && name == s.Name() {
					return nil, fmt.Errorf("named route path segment %q occurs twice", name)
				}
			}
			r.namedSegments = append(r.namedSegments, s)
		}
		r.segments = append(r.segments, s)
	}
	return r, nil
}

func (r *Route) Name() string {
	return r.name
}

func (r *Route) Fields(matchedValues []string) map[string]string {
	fields := make(map[string]string)
	for i, segment := range r.namedSegments {
		fields[segment.Name()] = matchedValues[i]
	}
	return fields
}

func (r *Route) Path(fields map[string]string) (string, error) {
	var (
		b     strings.Builder
		value string
		ok    bool
	)

	for _, s := range r.segments {
		if name := s.Name(); name != "" {
			_ = b.WriteByte('/')
			value, ok = fields[name]
			if !ok {
				return "", fmt.Errorf("field set for route %q does not contain field named %q", r, name)
			}
			_, _ = b.WriteString(value)
		} else {
			_, _ = b.WriteString(s.String())
		}
	}
	return b.String(), nil
}

func (r *Route) String() string {
	b := strings.Builder{}
	// _, _ = b.WriteString(r.name)
	// _ = b.WriteByte('=')
	if len(r.segments) == 0 {
		_, _ = b.WriteString("<nil>")
	}
	for _, s := range r.segments {
		_, _ = b.WriteString(s.String())
	}
	return b.String()
}
