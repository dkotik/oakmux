package oakmux

import (
	"fmt"
	"strings"
	"testing"

	_ "embed" // for test data
)

//go:embed test/data/strings.txt
var testStringsRaw string
var testStrings = strings.Split(strings.TrimSpace(testStringsRaw), "\n")

func BenchmarkNodeBranches(b *testing.B) {
	for i := 2; i < 25; i++ {
		searchKey := testStrings[i-1]
		var node *Node
		b.Run(fmt.Sprintf("branchList of %d items", i), func(b *testing.B) {
			var branches = make(branchList, 0, i)
			for _, key := range testStrings[:i] {
				branches = append(branches, keyedBranch{
					key:  key,
					node: &Node{},
				})
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				node = branches.Get(searchKey)
			}
		})
		if node == nil {
			b.Fatalf("search key %q was never found during  benchmark", searchKey)
		}

		b.Run(fmt.Sprintf("branchMap of %d items", i), func(b *testing.B) {
			var branches = make(branchMap)
			for _, key := range testStrings[:i] {
				branches[key] = &Node{}
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				node = branches.Get(searchKey)
			}
		})
		if node == nil {
			b.Fatalf("search key %q was never found during  benchmark", searchKey)
		}
	}
}
