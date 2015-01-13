package helpers

import "github.com/paulmach/go.geo/clustering/point_clustering"

// RemoveOutlierPointersByQuadkey will bucket all pointers by quad key (defined by the level)
// and remove the buckets with less than threshold pointers. The buckets become the resulting point_clustering.Clusters.
func RemoveOutlierPointersByQuadkey(pointers []point_clustering.Pointer, level, threshold int) []*point_clustering.Cluster {

	buckets := make(map[int64][]point_clustering.Pointer)
	for _, p := range pointers {
		key := p.CenterPoint().Quadkey(level)

		buckets[key] = append(buckets[key], p)
	}

	clusters := make([]*point_clustering.Cluster, 0, len(buckets))
	for _, b := range buckets {
		if len(b) >= threshold {
			clusters = append(clusters, point_clustering.NewCluster(b...))
		}
	}

	return clusters
}
