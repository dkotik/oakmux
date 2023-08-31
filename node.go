package oakmux

import (
	"fmt"
	"log"
	"strings"
)

// Node is the nesting routing tree component.
type Node struct {
	Leaf              *Route
	TrailingSlashLeaf *Route
	TerminalLeaf      *Route
	Branches          Branches
	DynamicBranch     *Node
}

type WalkFunc func(*Node) (ok bool, err error)

func (n *Node) MatchPath(path string) (route *Route, matches []string) {
	switch path {
	case "":
		return n.Leaf, nil
	case "/":
		return n.TrailingSlashLeaf, nil
	}

	segment, remainder, err := munchPath(path) // speed this up
	if err != nil {                            // TODO: drop this check
		log.Println("double slash error?", err)
		return nil, nil // double slash error possible only
	}
	// log.Println(segment, remainder)
	// spew.Dump(n)

	if n.Branches != nil {
		if branch := n.Branches.Get(segment); branch != nil {
			route, matches = branch.MatchPath(remainder)
			if route != nil {
				return route, matches
			}
		}
	}

	if n.DynamicBranch != nil {
		// TODO: pass in array with pre-initialized len instead?
		route, matches = n.DynamicBranch.MatchPath(remainder)
		if route != nil {
			return route, append([]string{segment}, matches...)
		}
	}

	if n.TerminalLeaf != nil {
		// if n.TerminalLeaf.segments[len(n.TerminalLeaf.segments)-1].Name() == "" {
		// 	return n.TerminalLeaf, nil // deal with the {...}
		// }
		return n.TerminalLeaf, []string{path[1:]}
	}
	return nil, nil
}

func (n *Node) Grow(route *Route, remaining []Segment) (err error) {
	if len(remaining) == 0 { // leaf
		if n.Leaf != nil {
			return fmt.Errorf("routes %q and %q overlap: %s resolves to the same static tree node as %s", n.Leaf.Name(), route.Name(), n.Leaf.String(), route.String())
		}
		n.Leaf = route
		return nil
	}
	current := remaining[0]
	switch current.Type() {
	case SegmentTypeTrailingSlash: // leaf
		if n.TrailingSlashLeaf != nil {
			return fmt.Errorf("routes %q and %q overlap: %s resolves to the same trailing slash tree node as %s", n.TrailingSlashLeaf.Name(), route.Name(), n.TrailingSlashLeaf.String(), route.String())
		}
		n.TrailingSlashLeaf = route
	case SegmentTypeTerminal: // leaf
		if n.TerminalLeaf != nil {
			return fmt.Errorf("routes %q and %q overlap: %s resolves to the same terminal tree node as %s", n.TerminalLeaf.Name(), route.Name(), n.TerminalLeaf.String(), route.String())
		}
		n.TerminalLeaf = route
	case SegmentTypeStatic: // branch
		if n.Branches == nil {
			n.Branches = make(branchList, 0, 1)
		}
		var node *Node
		// spew.Dump(current.Name())
		node, n.Branches = n.Branches.Grow(current.Name())
		// name := current.Name()
		// node := n.Branches.Get(name)
		// if node == nil {
		// 	node = &Node{}
		// 	n.Branches = n.Branches.Append(name, node)
		// }
		return node.Grow(route, remaining[1:])
	case SegmentTypeDynamic: // branch
		if n.DynamicBranch == nil {
			n.DynamicBranch = &Node{}
		}
		return n.DynamicBranch.Grow(route, remaining[1:])
	default:
		return fmt.Errorf("cannot grow tree using a segment %q of unknown type %q", current.Name(), current.Type())
	}
	return nil
}

func (n *Node) Walk(walkFn WalkFunc) (err error) {
	ok, err := walkFn(n)
	if err != nil || !ok {
		return
	}

	if n.Branches != nil {
		for _, key := range n.Branches.Keys() {
			if err = n.Branches.Get(key).Walk(walkFn); err != nil {
				return
			}
		}
	}

	if n.DynamicBranch != nil {
		return n.DynamicBranch.Walk(walkFn)
	}
	return
}

func (n *Node) String() string {
	b := &strings.Builder{}
	// b.WriteString("╟ ")
	if n.Leaf != nil {
		fmt.Fprintf(b, "[route:%s] ", n.Leaf.Name())
	}
	if n.TrailingSlashLeaf != nil {
		fmt.Fprintf(b, "[route/%s] ", n.TrailingSlashLeaf.Name())
	}
	if n.TerminalLeaf != nil {
		fmt.Fprintf(b, "[...%s]", n.TerminalLeaf.Name())
	}

	if n.Branches != nil {
		for _, branch := range n.Branches.Keys() {
			sub := n.Branches.Get(branch)
			b.WriteString("\n╚ ")
			fmt.Fprintf(b, "<%s>", branch)
			b.WriteString(strings.Replace(sub.String(), "\n", "\n    ", -1))
		}
	}

	if n.DynamicBranch != nil {
		b.WriteString("\n╚ <...>")
		b.WriteString(strings.Replace(n.DynamicBranch.String(), "\n", "\n    ", -1))
	}
	// fmt.Fprintf(b, "\n╟\n")
	// fmt.Fprintf(b, "╚ ╟\n")
	// Leaf              *Route
	// TrailingSlashLeaf *Route
	// TerminalLeaf      *Route
	// Branches          Branches
	// DynamicBranch     *Node
	return b.String()
}
