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
		r, err := NewRoute("test", testCase.Route)
		if err != nil {
			t.Fatal("cannot make route", testCase.Route, err)
		}
		if err = tree.Grow(r, r.segments); err != nil {
			t.Fatal("cannot grow tree node:", testCase.Route, err)
		}
	}
}

func TestNodeWalk(t *testing.T) {
	handler := newTestHandler(t)
	router, err := New(
		WithLimitlessRequestBytes(), // for type assertion
		WithRouteHandler(
			"firstRoute", "/test/[pattern]/yep/1/2/3/4",
			handler,
		),
		WithRouteHandler(
			"secondRoute", "/test/[wild]/[pattern1]/last/",
			handler,
		),
		WithRouteHandler(
			"thirdRoute", "/test/[pattern]/1/2/3/4",
			handler,
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	nodeCount := 0
	if err = router.(*mux).tree.Walk(func(n *Node) (ok bool, err error) {
		nodeCount++
		return true, nil
	}); err != nil {
		t.Fatal("failed to walk the Node tree:", err)
	}

	const expectedToVisit = 14
	if nodeCount != expectedToVisit {
		t.Fatalf("walk function was unable to visit every tree node: visited %d out of %d", nodeCount, expectedToVisit)
	}
}
