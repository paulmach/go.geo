package path_clustering

import "github.com/paulmach/go.geo"

// A Pather is the interface for something that can be path clustered.
type Pather interface {
	// RepresentativePath is kind of a weird name, but it's meant to not overlap
	// with any stuct attributes.
	RepresentativePath() *geo.Path
}

// A Cluster defines a cluster of paths.
type Cluster struct {
	Pathers []Pather
	indexes []int
}

// NewCluster creates a new Cluster properly.
func NewCluster(pathers ...Pather) *Cluster {
	c := &Cluster{
		Pathers: pathers,
	}

	c.indexes = make([]int, len(pathers))
	return c
}

// Merge merges the given cluster into the current cluster and returns. It basically
// just copies the pathers into itself.
func (c *Cluster) Merge(c2 *Cluster) {
	c.Pathers = append(c.Pathers, c2.Pathers...)
	c.indexes = append(c.indexes, c2.indexes...)
	return
}
