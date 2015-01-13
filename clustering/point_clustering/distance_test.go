package point_clustering

import "testing"

func TestCentroidDistance(t *testing.T) {
	// will not compile if interfaces not satisfied.
	var _ ClusterDistancer = CentroidDistance{}
	var _ ClusterDistancerFactory = CentroidDistance{}
}

func TestCentroidSquaredDistance(t *testing.T) {
	// will not compile if interfaces not satisfied.
	var _ ClusterDistancer = CentroidSquaredDistance{}
	var _ ClusterDistancerFactory = CentroidSquaredDistance{}
}

func TestCentroidGeoDistance(t *testing.T) {
	// will not compile if interfaces not satisfied.
	var _ ClusterDistancer = CentroidGeoDistance{}
	var _ ClusterDistancerFactory = CentroidGeoDistance{}
}
