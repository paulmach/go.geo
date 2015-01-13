package helpers

import (
	"github.com/paulmach/go.geo/clustering/path_clustering"
	"github.com/paulmach/go.geo/clustering/point_clustering"
)

// FilterSmallPointClusters will remove points clusters with less than or equal to the minPoints.
func FilterSmallPointClusters(clusters []*point_clustering.Cluster, minPoints int) []*point_clustering.Cluster {
	filtered := make([]*point_clustering.Cluster, 0, len(clusters))
	for _, c := range clusters {
		if len(c.Pointers) >= minPoints {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

// FilterSmallPathClusters will remove path clusters with less than or equal to the minPaths.
func FilterSmallPathClusters(clusters []*path_clustering.Cluster, minPaths int) []*path_clustering.Cluster {
	filtered := make([]*path_clustering.Cluster, 0, len(clusters))
	for _, c := range clusters {
		if len(c.Pathers) >= minPaths {
			filtered = append(filtered, c)
		}
	}

	return filtered
}
