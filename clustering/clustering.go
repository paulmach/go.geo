package clustering

import (
	"math"

	"github.com/paulmach/go.geo"
)

// ClusterPointers will take a set of Pointers and cluster them using
// the distancer and threshold.
func ClusterPointers(pointers []Pointer, distancer ClusterDistancer, threshold float64) []*Cluster {
	clusters := make([]*Cluster, 0, len(pointers))
	for _, p := range pointers {
		clusters = append(clusters, NewCluster(p))
	}

	return cluster(clusters, distancer, threshold)
}

// ClusterClusters can be used if you've already created cluster objects
// using a prefilterer of something else. Original clusters will be copied
// so the original set will be unchanged.
func ClusterClusters(clusters []*Cluster, distancer ClusterDistancer, threshold float64) []*Cluster {
	copiedClusters := make([]*Cluster, len(clusters), len(clusters))
	copy(copiedClusters, clusters)

	return cluster(copiedClusters, distancer, threshold)
}

// cluster will modify the passed in clusters, centroid and list of pointers,
// so a copy must have been made before reaching this function.
func cluster(clusters []*Cluster, distancer ClusterDistancer, threshold float64) []*Cluster {
	if len(clusters) < 2 {
		return clusters
	}

	count := 0
	for _, cluster := range clusters {
		count += len(cluster.Pointers)
	}

	clusters, found := clusterClusters(
		clusters,
		// Default intialization, TODO: better bucketing/prefiltering will greatly increase performance.
		initClusterDistances(clusters, distancer, threshold),
		distancer,
		threshold,
	)

	result := make([]*Cluster, 0, found)
	for _, cluster := range clusters {
		if cluster != nil {
			result = append(result, cluster)
		}
	}

	return result
}

// ClusterGeoPointers will take a set of Pointers and cluster them.
// It will project the points using mercator, scale the threshold, cluster, and project back.
// Performace is about 40% than simply using a geo distancer.
// This may not make sense for all geo datasets.
func ClusterGeoPointers(pointers []Pointer, threshold float64) []*Cluster {
	clusters := make([]*Cluster, 0, len(pointers))
	for _, p := range pointers {
		clusters = append(clusters, NewCluster(p))
	}

	if len(clusters) < 2 {
		return clusters
	}

	// performs the actual clustering
	return geocluster(clusters, threshold)
}

// ClusterGeoClusters can be used if you've already created clusters objects
// using a prefilterer of something else.
func ClusterGeoClusters(clusters []*Cluster, threshold float64) []*Cluster {
	if len(clusters) < 2 {
		return clusters
	}

	copiedClusters := make([]*Cluster, len(clusters), len(clusters))
	for i, cluster := range clusters {
		copiedClusters[i] = NewClusterWithCentroid(cluster.Centroid, cluster.Pointers...)
	}

	return geocluster(copiedClusters, threshold)
}

// will modify the passed in clusters, centroid and list of pathers,
// so a copy must have been made before reaching this function.
func geocluster(clusters []*Cluster, threshold float64) []*Cluster {
	if len(clusters) < 2 {
		return clusters
	}

	bound := geo.NewBoundFromPoints(clusters[0].Centroid, clusters[0].Centroid)
	for _, cluster := range clusters {
		bound.Extend(cluster.Centroid)
		geo.Mercator.Project(cluster.Centroid)
	}

	factor := geo.MercatorScaleFactor(bound.Center().Lat())
	scaledThreshold := threshold * threshold * factor * factor

	clusteredClusters, found := clusterClusters(
		clusters,
		// Default intialization, TODO: better bucketing/prefiltering will greatly increase performance.
		// can use the bound above to help with this.
		initClusterDistances(clusters, CentroidSquaredDistance{}, scaledThreshold),
		CentroidSquaredDistance{},
		scaledThreshold,
	)

	result := make([]*Cluster, 0, found)
	for _, cluster := range clusteredClusters {
		if cluster != nil {
			geo.Mercator.Inverse(cluster.Centroid)
			result = append(result, cluster)
		}
	}

	return result
}

func initClusterDistances(
	clusters []*Cluster,
	distancer ClusterDistancer,
	threshold float64,
) []*distanceSet {

	// initialize distances
	distances := make([]*distanceSet, len(clusters))

	for i := 0; i < len(clusters); i++ {

		if distances[i] == nil {
			distances[i] = newDistanceSet()
		}
		distances[i].Set(i, math.MaxInt32)

		for j := i + 1; j < len(clusters); j++ {
			// TODO: better filtering here we don't have a literal n^2 situation.
			dist := distancer.ClusterDistance(clusters[i], clusters[j])
			if dist < 5*threshold {
				distances[i].Set(j, dist)

				if distances[j] == nil {
					distances[j] = newDistanceSet()
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
	distanceSets []*distanceSet,
	distancer ClusterDistancer,
	threshold float64,
) ([]*Cluster, int) {

	s := &state{
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
		clusters[lower] = CombineClusters(clusters[lower], clusters[higher])
		s.ResetDistances(lower, higher)
		clusters[higher] = nil

		removed++
	}

	return clusters, len(clusters) - removed
}
