package oakmux

const optimalMimimumBranchMapSize = 8

// Branches abstracts either map or list implementation for child [Node]s for performance. When there are more than [optimalMimimumBranchMapSize], the map implementation is prefered for faster look up.
//
// TODO: John Amsterdam's implementation switched to generics instead of
// hybrid. Would it be faster than asserting the interface? Should not be.
type Branches interface {
	Get(string) *Node
	Grow(string) (*Node, Branches)
	Keys() []string
}

var _ Branches = (branchList)(nil) // ensure interface satisfaction
var _ Branches = (branchMap)(nil)  // ensure interface satisfaction

type keyedBranch struct {
	key  string
	node *Node
}

type branchList []keyedBranch

func (l branchList) Get(key string) *Node {
	for _, c := range l {
		if c.key == key {
			return c.node
		}
	}
	return nil
}

func (l branchList) Grow(key string) (*Node, Branches) {
	node := l.Get(key)
	if node != nil { // no need to grow
		return node, l
	}
	if len(l) >= optimalMimimumBranchMapSize { // need to become a map
		m := make(branchMap)
		for _, c := range l {
			m[c.key] = c.node
		}
		return m.Grow(key)
	}
	node = &Node{}
	return node, append(l, keyedBranch{key, node})
}

func (l branchList) Keys() []string {
	// TODO: if keys are used frequently, they can be in a separate slice.
	// "golang.org/x/exp/slices" may offer some utility functions that speed
	// up [branchList] implementation.
	//
	// TODO: What does slices.Compact([]string) do? frees up tail memory?
	keys := make([]string, len(l))
	for i, c := range l {
		keys[i] = c.key
	}
	return keys
}

type branchMap map[string]*Node

func (m branchMap) Get(key string) (node *Node) {
	node, ok := m[key]
	if ok {
		return node
	}
	return nil
}

func (m branchMap) Grow(key string) (*Node, Branches) {
	node, ok := m[key]
	if ok {
		return node, m
	}
	node = &Node{}
	m[key] = node
	return node, m
}

func (m branchMap) Keys() []string {
	// TODO: "golang.org/x/exp/maps" has maps.Keys(ms) method, may be faster.
	keys := make([]string, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
