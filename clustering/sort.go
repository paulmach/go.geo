package clustering

import (
	"sort"
)

// Sortable implements the sorting interface
// allowing for sorting.
type Sortable []*Cluster

// Sort will sort with set.
// Usage: clustering.Sortable(clusters).Sort()
func (s Sortable) Sort() {
	sort.Sort(s)
}

// Len returns the length of the sortable cluster.
// This is to implement the sort interface.
func (s Sortable) Len() int {
	return len(s)
}

// Less returns truee if i > j, so bigger will be first.
// This is to implement the sort interface.
func (s Sortable) Less(i, j int) bool {
	return len(s[i].Pointers) > len(s[j].Pointers)

}

// Swap interchanges two elements.
// This is to implement the sort interface.
func (s Sortable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
