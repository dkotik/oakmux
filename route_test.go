package oakmux

import (
	"errors"
	"testing"
)

func TestSuccessfulRouteCreation(t *testing.T) {
	cases := []string{
		"/1/2/3/4/5/6",
		"/good/routes",
		"/a/b/c/",
	}
	for _, testCase := range cases {
		r, err := NewRoute("test", testCase)
		if err != nil {
			t.Fatal("cannot make route", testCase, err)
		}
		reconstitute := r.String()
		if reconstitute != testCase {
			t.Fatal("reconstituted path does not match its original:",
				reconstitute, testCase)
		}
	}
}

func TestForDoubleSlashes(t *testing.T) {
	cases := []string{
		"//1/2/3/4/5/6",
		"/good//routes",
		"/a/b/c//",
		"///a/b/c/",
		"/a/b//c/",
		"/a/b//c////",
	}
	for _, testCase := range cases {
		_, err := NewRoute("test", testCase)
		if !errors.Is(err, ErrDoubleSlash) {
			t.Fatalf("expected double slashes error for %q, but got %q error instead", testCase, err)
		}
	}
}
