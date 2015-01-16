package helpers

import "github.com/paulmach/go.geo/clustering"

// RemoveOutlierPointersByQuadkey will bucket all pointers by quad key (defined by the level)
// and remove the buckets with less than threshold pointers. The buckets become the resulting point_clustering.Clusters.
func RemoveOutlierPointersByQuadkey(pointers []clustering.Pointer, level, threshold int) []*clustering.Cluster {

	buckets := make(map[int64][]clustering.Pointer)
	for _, p := range pointers {
		key := p.CenterPoint().Quadkey(level)

		buckets[key] = append(buckets[key], p)
	}

	clusters := make([]*clustering.Cluster, 0, len(buckets))
	for _, b := range buckets {
		if len(b) >= threshold {
			clusters = append(clusters, clustering.NewCluster(b...))
		}
	}

	return clusters
}
