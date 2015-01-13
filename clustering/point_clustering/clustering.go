package point_clustering

import (
	"math"

	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/clustering/shared"
)

// Clustering defines parameters for the point clustering algorithm.
type Clustering struct {
	Threshold        float64
	DistancerFactory ClusterDistancerFactory
}

// New creates a new point clustering config object.
func New(threshold float64, factory ClusterDistancerFactory) *Clustering {
	return &Clustering{
		Threshold:        threshold,
		DistancerFactory: factory,
	}
}

// Cluster will take a set of Pathers and cluster them using the distance threshold
// and other paramters from the Clustering struct.
func (c *Clustering) Cluster(pointers []Pointer) []*Cluster {
	clusters := make([]*Cluster, 0, len(pointers))
	for _, p := range pointers {
		clusters = append(clusters, NewCluster(p))
	}

	// performs the actual clustering
	return c.cluster(clusters)
}

// ClusterClusters can be used if you've already created cluster objects
// using a prefilterer of something else.
func (c *Clustering) ClusterClusters(clusters []*Cluster) []*Cluster {
	copiedClusters := make([]*Cluster, len(clusters), len(clusters))
	for i, cluster := range clusters {
		copiedClusters[i] = NewClusterWithCentroid(cluster.Centroid, cluster.Pointers...)
	}

	return c.cluster(copiedClusters)
}

// cluster will modify the passed in clusters, centroid and list of pointers,
// so a copy must have been made before reaching this function.
func (c *Clustering) cluster(clusters []*Cluster) []*Cluster {
	if len(clusters) < 2 {
		return clusters
	}

	count := 0
	for _, cluster := range clusters {
		count += len(cluster.Pointers)
	}

	distancer := c.DistancerFactory.ClusterDistancer(len(clusters), count)
	clusters, found := clusterClusters(
		clusters,
		// Default intialization, TODO: better bucketing/prefiltering will greatly increase performance.
		initializeClusterDistances(clusters, distancer, c.Threshold),
		distancer,
		c.Threshold,
	)

	result := make([]*Cluster, 0, found)
	for _, cluster := range clusters {
		if cluster != nil {
			result = append(result, cluster)
		}
	}

	return result
}

// GeoProjectedClustering defines parameters for the clustering algorithm.
// This clustering will project all the points to a mercator projection (EPSG:3857)
// and use a euclidean distance function, scaled appropriately.
type GeoProjectedClustering struct {
	Threshold float64
}

// NewPointClusteringGeoProjected creates a new point clustering config object.
// This clustering must be used only for geo (lng/lat) points.
func NewGeoProjectedClustering(threshold float64) *GeoProjectedClustering {
	return &GeoProjectedClustering{
		Threshold: threshold,
	}
}

// Cluster will take a set of Pointers and cluster them using the distance threshold
// and other paramters from the PointClusteringGeoProjected struct.
func (c *GeoProjectedClustering) Cluster(pointers []Pointer) []*Cluster {
	clusters := make([]*Cluster, 0, len(pointers))
	for _, p := range pointers {
		clusters = append(clusters, NewCluster(p))
	}

	if len(clusters) < 2 {
		return clusters
	}

	// performs the actual clustering
	return c.cluster(clusters)
}

// ClusterClusters can be used if you've already created clusters objects
// using a prefilterer of something else.
func (c *GeoProjectedClustering) ClusterClusters(clusters []*Cluster) []*Cluster {
	if len(clusters) < 2 {
		return clusters
	}

	copiedClusters := make([]*Cluster, len(clusters), len(clusters))
	for i, cluster := range clusters {
		copiedClusters[i] = NewClusterWithCentroid(cluster.Centroid, cluster.Pointers...)
	}

	return c.cluster(copiedClusters)
}

// cluster will modify the passed in clusters, centroid and list of pathers,
// so a copy must have been made before reaching this function.
func (c *GeoProjectedClustering) cluster(clusters []*Cluster) []*Cluster {
	if len(clusters) < 2 {
		return clusters
	}

	bound := geo.NewBoundFromPoints(clusters[0].Centroid, clusters[1].Centroid)
	for _, cluster := range clusters {
		bound.Extend(cluster.Centroid)
		geo.Mercator.Project(cluster.Centroid)
	}

	factor := geo.MercatorScaleFactor(bound.Center().Lat())
	threshold := c.Threshold * c.Threshold * factor * factor

	clusteredClusters, found := clusterClusters(
		clusters,
		// Default intialization, TODO: better bucketing/prefiltering will greatly increase performance.
		// can use the bound above to help with this.
		initializeClusterDistances(clusters, CentroidSquaredDistance{}, threshold),
		CentroidSquaredDistance{},
		threshold,
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

func initializeClusterDistances(clusters []*Cluster, distancer ClusterDistancer, threshold float64) []*shared.DistanceSet {
	// initialize distances
	distances := make([]*shared.DistanceSet, len(clusters))
	for i := 0; i < len(clusters); i++ {

		if distances[i] == nil {
			distances[i] = shared.NewDistanceSet()
		}
		distances[i].Set(i, math.MaxInt32)

		for j := i + 1; j < len(clusters); j++ {
			// TODO: better filtering here we don't have a literal n^2 situation.
			dist := distancer.ClusterDistance(clusters[i], clusters[j])
			if dist < 5*threshold {
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
