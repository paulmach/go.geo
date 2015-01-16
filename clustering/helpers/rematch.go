package helpers

import (
	"math"

	"github.com/paulmach/go.geo/clustering"
)

// RematchPointersToClusters will take a set of pointers and map them to the closest cluster.
// Basically creates a new cluster from that one point and does the ClusterDistance between them.
// Will return a new list.
func RematchPointersToClusters(
	clusters []*clustering.Cluster,
	pointers []clustering.Pointer,
	distancer clustering.ClusterDistancer,
	threshold float64,
) []*clustering.Cluster {
	if len(clusters) == 0 {
		return []*clustering.Cluster{}
	}

	newClusters := make([]*clustering.Cluster, 0, len(clusters))

	// clear the current members
	for _, c := range clusters {
		newClusters = append(newClusters, clustering.NewClusterWithCentroid(c.Centroid))
	}

	// remap all the groupers to these new groups
	for _, pointer := range pointers {
		minDist := math.MaxFloat64
		index := 0

		pointerCluster := clustering.NewCluster(pointer)

		// find the closest group
		for i, c := range newClusters {
			if d := distancer.ClusterDistance(c, pointerCluster); d < minDist {
				minDist = d
				index = i
			}
		}

		if minDist < threshold {
			// leaves the center as found by the previous clustering
			newClusters[index].Pointers = append(newClusters[index].Pointers, pointer)
		}
	}

	return newClusters
}
