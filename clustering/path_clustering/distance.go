package path_clustering

import (
	"math"

	"github.com/paulmach/go.geo/clustering/shared"
)

// PatherDistancer defines exactly how the distance
// between two paths is computed.
type PatherDistancer interface {
	PatherDistance(p1, p2 Pather) float64
}

// PatherDistancerFactory builds PatherDistancers. Since the distancer
// may cache values we need to build a new one for every run using this interface.
type PatherDistancerFactory interface {
	PatherDistancer(pathCount int) PatherDistancer
}

// ClusterDistancer defines the distance between path clusters.
// Since the distance between two sets can be defined different ways,
// one just needs to implement this interface to get the desired behavior.
// Objects implementing this interface should do some sort of PatherDistance caching.
type ClusterDistancer interface {
	ClusterDistance(pc1, pc2 *Cluster) float64
}

// PathClusterDistancerFactory builds PathClusterDistancers.
// Since these distancers can also hold state, a different one must be used
// for each clustering operations. That new distancer is built using this interface.
type ClusterDistancerFactory interface {
	// PathClusterDistancer clusterCount and pathCount paramters can be used
	// to initialize caching data structures.
	ClusterDistancer(clusterCount, pathCount int) ClusterDistancer
}

// SingleLinkageDistancerFactory builds SingleLinkageDistancers.
type SingleLinkageDistancerFactory struct {
	DistancerFactory PatherDistancerFactory
}

// NewSingleLinkageDistancerFactory will create an object to build this type of distancer.
func NewSingleLinkageDistancerFactory(factory PatherDistancerFactory) *SingleLinkageDistancerFactory {
	return &SingleLinkageDistancerFactory{
		DistancerFactory: factory,
	}
}

// PathClusterDistancer returns a newly minted PathClusterSingleLinkageDistancer.
func (f *SingleLinkageDistancerFactory) ClusterDistancer(clusterCount, pathCount int) ClusterDistancer {
	return NewSingleLinkageDistancer(
		f.DistancerFactory.PatherDistancer(pathCount),
		clusterCount,
		pathCount,
	)
}

// SingleLinkageDistancer computes cluster distance by returning the minimum
// of all possible links between the clusters.
type SingleLinkageDistancer struct {
	cacher    shared.Cacher
	Distancer PatherDistancer
}

// NewPathClusterSingleLinkageDistancer creates a new PathClusterSingleLinkageDistancer.
// The caching object type is chosen based on the count parameters.
func NewSingleLinkageDistancer(distancer PatherDistancer, clusterCount, pathCount int) *SingleLinkageDistancer {
	d := &SingleLinkageDistancer{
		Distancer: distancer,
	}

	// 4000 == 16,000,000 elements, 128,000,000 megs.
	if pathCount <= 4000 {
		d.cacher = shared.NewArrayCache(pathCount)
	} else {
		d.cacher = shared.NewMapCache(pathCount)
	}

	return d
}

// PathClusterDistance computes the distance between the clusters using single linkage.
// ie. it takes the minimum of the distance between all possible links.
func (d *SingleLinkageDistancer) ClusterDistance(c1, c2 *Cluster) float64 {
	min := math.MaxFloat64
	for i, p1 := range c1.Pathers {
		for j, p2 := range c2.Pathers {
			dist := d.cacher.Get(c1.indexes[i], c2.indexes[j])
			if dist < 0 {
				dist = d.Distancer.PatherDistance(p1, p2)
				d.cacher.Set(c1.indexes[i], c2.indexes[j], dist)
			}

			min = math.Min(dist, min)
		}
	}

	return min
}

// CompleteLinkageDistancerFactory builds CompleteLinkageDistancers.
type CompleteLinkageDistancerFactory struct {
	DistancerFactory PatherDistancerFactory
}

// NewPathClusterCompleteLinkageDistancerFactory will create an object to build this type of distancer.
func NewCompleteLinkageDistancerFactory(factory PatherDistancerFactory) *CompleteLinkageDistancerFactory {
	return &CompleteLinkageDistancerFactory{
		DistancerFactory: factory,
	}
}

// PathClusterDistancer returns a newly minted CompleteLinkageDistancers.
func (f *CompleteLinkageDistancerFactory) ClusterDistancer(clusterCount, pathCount int) ClusterDistancer {
	return NewCompleteLinkageDistancer(
		f.DistancerFactory.PatherDistancer(pathCount),
		clusterCount,
		pathCount,
	)
}

// CompleteLinkageDistancers computes cluster distance by returning the maximum
// of all possible links between the clusters.
type CompleteLinkageDistancer struct {
	cacher    shared.Cacher
	Distancer PatherDistancer
}

// NewCompleteLinkageDistancers creates a new PathClusterCompleteLinkageDistancer.
// The caching object type is chosen based on the count parameters.
func NewCompleteLinkageDistancer(distancer PatherDistancer, clusterCount, pathCount int) *CompleteLinkageDistancer {
	d := &CompleteLinkageDistancer{
		Distancer: distancer,
	}

	// 4000 == 16,000,000 elements, 128,000,000 megs.
	if pathCount <= 4000 {
		d.cacher = shared.NewArrayCache(pathCount)
	} else {
		d.cacher = shared.NewMapCache(pathCount)
	}

	return d
}

// PathClusterDistance computes the distance between the clusters using complete linkage.
// ie. it takes the minimum of the distance between all possible links.
func (d *CompleteLinkageDistancer) ClusterDistance(c1, c2 *Cluster) float64 {
	max := 0.0
	for i, p1 := range c1.Pathers {
		for j, p2 := range c2.Pathers {
			dist := d.cacher.Get(c1.indexes[i], c2.indexes[j])
			if dist < 0 {
				dist = d.Distancer.PatherDistance(p1, p2)
				d.cacher.Set(c1.indexes[i], c2.indexes[j], dist)
			}

			max = math.Max(dist, max)
		}
	}

	return max
}
