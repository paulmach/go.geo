package helpers

import "github.com/paulmach/go.geo/clustering"

// FilterSmallClusters will remove points clusters with less than or equal to the minPoints.
func FilterSmallClusters(clusters []*clustering.Cluster, minPoints int) []*clustering.Cluster {
	filtered := make([]*clustering.Cluster, 0, len(clusters))
	for _, c := range clusters {
		if len(c.Pointers) >= minPoints {
			filtered = append(filtered, c)
		}
	}

	return filtered
}
