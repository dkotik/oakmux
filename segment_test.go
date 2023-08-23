package oakmux

import (
	"bytes"
	"testing"
)

func TestSegmentCreation(t *testing.T) {
	t.Run("SegmentTypeStatic", func(t *testing.T) {
		cases := []struct {
			Path   string
			Result staticSegment
		}{
			{Path: "/var", Result: staticSegment("var")},
			{Path: "/crazy", Result: staticSegment("crazy")},
			{Path: "something", Result: staticSegment("something")},
		}

		for _, testCase := range cases {
			t.Run(testCase.Path, func(t *testing.T) {
				segment, err := NewSegment(testCase.Path)
				if err != nil {
					t.Fatal(err)
				}
				cast, ok := segment.(staticSegment)
				if !ok {
					t.Fatalf("types %T and %T do not match", cast, segment)
				}
				if bytes.Compare(cast, testCase.Result) != 0 {
					t.Fatalf("%q != %q", cast, testCase.Result)
				}
			})
		}
	})
	t.Run("SegmentTypeDynamic", func(t *testing.T) {
		cases := []struct {
			Path   string
			Result dynamicSegment
		}{
			{Path: "/{var}", Result: dynamicSegment("var")},
			{Path: "/{}", Result: dynamicSegment("")},
			{Path: "{something}", Result: dynamicSegment("something")},
		}

		for _, testCase := range cases {
			t.Run(testCase.Path, func(t *testing.T) {
				segment, err := NewSegment(testCase.Path)
				if err != nil {
					t.Fatal(err)
				}
				cast, ok := segment.(dynamicSegment)
				if !ok {
					t.Fatalf("types %T and %T do not match", cast, segment)
				}
				if cast != testCase.Result {
					t.Fatalf("%q != %q", cast, testCase.Result)
				}
			})
		}
	})
	t.Run("SegmentTypeTerminal", func(t *testing.T) {
		cases := []struct {
			Path   string
			Result terminalSegment
		}{
			{Path: "/{...var}", Result: terminalSegment("var")},
			{Path: "/{...}", Result: terminalSegment("")},
			{Path: "{...something}", Result: terminalSegment("something")},
		}

		for _, testCase := range cases {
			t.Run(testCase.Path, func(t *testing.T) {
				segment, err := NewSegment(testCase.Path)
				if err != nil {
					t.Fatal(err)
				}
				cast, ok := segment.(terminalSegment)
				if !ok {
					t.Fatalf("types %T and %T do not match", cast, segment)
				}
				if cast != testCase.Result {
					t.Fatalf("%q != %q", cast, testCase.Result)
				}
			})
		}
	})
	t.Run("SegmentTypeTrailingSlash", func(t *testing.T) {
		segment, err := NewSegment("/")
		if err != nil {
			t.Fatal(err)
		}
		cast, ok := segment.(trailingSlashSegment)
		if !ok {
			t.Fatalf("types %T and %T do not match", cast, segment)
		}
	})
}
