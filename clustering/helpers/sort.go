package helpers

import (
	"sort"

	"github.com/paulmach/go.geo/clustering"
)

// SortableClusters can be sorted so that clusters
// with more elements are first.
type SortableClusters []*clustering.Cluster

// Sort will sort with set.
// Usage: helpers.SortableClusters(clusters).Sort()
func (s SortableClusters) Sort() {
	sort.Sort(s)
}

// Len returns the length of the sortable cluster.
// This is to implement the sort interface.
func (s SortableClusters) Len() int {
	return len(s)
}

// Less returns truee if i > j, so bigger will be first.
// This is to implement the sort interface.
func (s SortableClusters) Less(i, j int) bool {
	return len(s[i].Pointers) > len(s[j].Pointers)

}

// Swap interchanges two elements.
// This is to implement the sort interface.
func (s SortableClusters) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
