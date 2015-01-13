package helpers

import (
	"math"

	"github.com/paulmach/go.geo/clustering/point_clustering"
)

// RematchPointersToClusters will take a set of pointers and map them to the closest cluster.
// Basically creates a new cluster from that one point and does the ClusterDistance between them.
// Will return a new list.
func RematchPointersToClusters(
	clusters []*point_clustering.Cluster,
	pointers []point_clustering.Pointer,
	distancer point_clustering.ClusterDistancer,
	threshold float64,
) []*point_clustering.Cluster {
	if len(clusters) == 0 {
		return []*point_clustering.Cluster{}
	}

	newClusters := make([]*point_clustering.Cluster, 0, len(clusters))

	// clear the current members
	for _, c := range clusters {
		newClusters = append(newClusters, point_clustering.NewClusterWithCentroid(c.Centroid))
	}

	// remap all the groupers to these new groups
	for _, pointer := range pointers {
		minDist := math.MaxFloat64
		index := 0

		pointerCluster := point_clustering.NewCluster(pointer)

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
