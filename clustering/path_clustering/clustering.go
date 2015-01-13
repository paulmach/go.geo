package path_clustering

import (
	"math"

	"github.com/paulmach/go.geo/clustering/shared"
)

// lustering defines parameters for the path clustering algorithm.
type Clustering struct {
	Threshold        float64
	DistancerFactory ClusterDistancerFactory
}

// New creates a new path clustering config.
func New(threshold float64, factory ClusterDistancerFactory) *Clustering {
	return &Clustering{
		Threshold:        threshold,
		DistancerFactory: factory,
	}
}

// Cluster will take a set of Pathers and cluster them using the distance threshold
// and Distancer from the Clustering struct.
func (c *Clustering) Cluster(pathers []Pather) []*Cluster {
	clusters := make([]*Cluster, 0, len(pathers))
	for _, p := range pathers {
		clusters = append(clusters, NewCluster(p))
	}

	// performs the actual clustering
	return c.cluster(clusters)
}

// ClusterClusters can be used if you've already created Cluster objects
// using a prefilterer of something else.
func (c *Clustering) ClusterClusters(clusters []*Cluster) []*Cluster {
	copiedClusters := make([]*Cluster, len(clusters), len(clusters))
	for i, cluster := range clusters {
		copiedClusters[i] = NewCluster(cluster.Pathers...)
	}

	return c.cluster(copiedClusters)
}

// cluster will modify the passed in clusers and list of pathers,
// so a copy must have been made before reaching this function.
func (c *Clustering) cluster(clusters []*Cluster) []*Cluster {
	// set indexes
	index := 0
	for _, cluster := range clusters {
		for i := range cluster.Pathers {
			cluster.indexes[i] = index
			index++
		}
	}

	distancer := c.DistancerFactory.ClusterDistancer(len(clusters), index)
	clusteredClusters, found := clusterClusters(
		clusters,
		initializeClusterDistances(clusters, distancer, c.Threshold),
		distancer,
		c.Threshold,
	)

	// remove nil values from result
	result := make([]*Cluster, 0, found)
	for _, c := range clusteredClusters {
		if c != nil {
			result = append(result, c)
		}
	}

	return result
}

func initializeClusterDistances(
	clusters []*Cluster,
	distancer ClusterDistancer,
	threshold float64,
) []*shared.DistanceSet {

	// initialize distances
	distances := make([]*shared.DistanceSet, len(clusters))
	for i := 0; i < len(clusters); i++ {
		if clusters[i] == nil {
			continue
		}

		if distances[i] == nil {
			distances[i] = shared.NewDistanceSet()
		}
		distances[i].Set(i, math.MaxInt32)

		for j := i + 1; j < len(clusters); j++ {
			if clusters[j] == nil {
				continue
			}

			// TODO: better filtering here we don't have a literal n^2 situation.
			dist := distancer.ClusterDistance(clusters[i], clusters[j])
			if dist < 10*threshold {
				distances[i].Set(j, dist)

				if distances[j] == nil {
					distances[j] = shared.NewDistanceSet()
				}
				distances[j].Set(i, dist)
			} else {
				// greater than a big threshold, so pass
			}
		}
	}

	return distances
}

func clusterClusters(
	clusters []*Cluster,
	distanceSets []*shared.DistanceSet,
	distancer ClusterDistancer,
	threshold float64,
) ([]*Cluster, int) {

	s := &shared.State{
		Distances: distanceSets,
		DistanceFunc: func(a, b int) float64 {
			return distancer.ClusterDistance(clusters[a], clusters[b])
		},
	}

	// successively merge
	removed := 0
	for len(clusters)-removed > 1 {
		lower, higher, dist := s.MinDistance()
		if dist > threshold {
			break
		}

		// merge these two
		clusters[lower].Merge(clusters[higher])
		s.ResetDistances(lower, higher)
		clusters[higher] = nil

		removed++
	}

	return clusters, len(clusters) - removed
}
