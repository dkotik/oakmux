package oakmux

import (
	"fmt"
	"strings"
)

type SegmentType uint8

const (
	SegmentTypeUnknown = iota
	SegmentTypeStatic
	SegmentTypeDynamic
	SegmentTypeTerminal
	SegmentTypeTrailingSlash
)

func (s SegmentType) String() string {
	switch s {
	case SegmentTypeStatic:
		return "static"
	case SegmentTypeDynamic:
		return "dynamic"
	case SegmentTypeTerminal:
		return "terminal"
	case SegmentTypeTrailingSlash:
		return "trailing slash"
	default:
		return "unknown"
	}
}

type Segment interface {
	Name() string
	Type() SegmentType
	Match(path string) (value, remainder string, ok bool)
	String() string
}

// NewSegment converts a string to a [Segment] definition.
func NewSegment(segmentDefinition string) (Segment, error) {
	switch segmentDefinition {
	case "":
		return nil, fmt.Errorf("empty path segment")
	case "/":
		return trailingSlashSegment{}, nil
	}

	if segmentDefinition[0] == '/' {
		segmentDefinition = segmentDefinition[1:]
	}
	if segmentDefinition[0] == '[' {
		tail := len(segmentDefinition) - 1
		if segmentDefinition[tail] != ']' {
			return nil, fmt.Errorf("dynamic path segment definition %q is missing a closing curly brace", segmentDefinition)
		}
		segmentDefinition = segmentDefinition[1:tail] // cut off {}
		if strings.HasPrefix(segmentDefinition, "...") && tail > 3 {
			return terminalSegment(segmentDefinition[3:]), nil
		}
		return dynamicSegment(segmentDefinition), nil
	}
	return staticSegment(segmentDefinition), nil
}

type staticSegment []byte

func (s staticSegment) Name() string {
	return string(s)
}

func (s staticSegment) Type() SegmentType {
	return SegmentTypeStatic
}

func (s staticSegment) Match(path string) (string, string, bool) {
	lengthSegment, lengthPath := len(s), len(path)
	if lengthSegment <= lengthPath {
		// presupposes that s[0] == '/'
		for i, c := range s {
			if path[i] != c {
				break
			}
		}
		if lengthSegment == lengthPath || path[lengthSegment] == '/' {
			return string(s), path[:lengthSegment], true
		}
	}
	return "", "", false
}

func (s staticSegment) String() string {
	return "/" + string(s)
}

type dynamicSegment string

func (d dynamicSegment) Name() string {
	return string(d)
}

func (d dynamicSegment) Type() SegmentType {
	return SegmentTypeDynamic
}

func (d dynamicSegment) Match(path string) (string, string, bool) {
	// if path == "" || path == "." {
	// 	return "", "", true
	// }
	if path[0] == '/' {
		path = path[1:] // ignore leading slash
	}

	var (
		c byte
		i int
	)
	for i, c = range []byte(path) {
		if c == '/' {
			break
		}
	}
	return path[:i], path[i:], true
}

func (d dynamicSegment) String() string {
	return "/[" + string(d) + "]"
}

type terminalSegment string

func (t terminalSegment) Name() string {
	return string(t)
}

func (t terminalSegment) Type() SegmentType {
	return SegmentTypeTerminal
}

func (t terminalSegment) Match(path string) (string, string, bool) {
	if path[0] == '/' {
		path = path[1:] // ignore leading slash
	}
	tail := len(path) - 1
	if tail == -1 { // path == ""
		return "", "", false // no content
	}
	if path[tail] == '/' {
		if tail >= 1 && path[tail-1] == '/' {
			return "", "", false // multi-slash tail, like "...//"
		}
		return path[:tail], path[tail:], true // trailing slash
	}
	return path, "", true
}

func (t terminalSegment) String() string {
	return "/[..." + string(t) + "]"
}

type trailingSlashSegment struct{}

func (t trailingSlashSegment) Name() string {
	return ""
}

func (t trailingSlashSegment) Type() SegmentType {
	return SegmentTypeTrailingSlash
}

func (t trailingSlashSegment) Match(path string) (string, string, bool) {
	return "", "", path == "/"
}

func (t trailingSlashSegment) String() string {
	return "/"
}
