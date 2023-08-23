package oakmux

import "testing"

func TestTreeCreation(t *testing.T) {
	cases := []struct {
		Route    string
		Matching []string
		Failing  []string
	}{
		{
			Route:    "/1/2",
			Matching: []string{"/1/2"},
			Failing:  []string{"/1/22"},
		},
	}
	// "/1/2/3/4/5/6",
	// "/good/routes",
	// "/a/b/c/",

	tree := &Node{}
	for _, testCase := range cases {
		r, err := NewRoute(testCase.Route)
		if err != nil {
			t.Fatal("cannot make route", testCase.Route, err)
		}
		if err = tree.Grow(r, r.segments); err != nil {
			t.Fatal("cannot grow tree node:", testCase.Route, err)
		}
	}
}
