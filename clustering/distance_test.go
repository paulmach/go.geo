package clustering

import "testing"

func TestCentroidDistance(t *testing.T) {
	// will not compile if interfaces not satisfied.
	var _ ClusterDistancer = CentroidDistance{}
}

func TestCentroidSquaredDistance(t *testing.T) {
	// will not compile if interfaces not satisfied.
	var _ ClusterDistancer = CentroidSquaredDistance{}
}

func TestCentroidGeoDistance(t *testing.T) {
	// will not compile if interfaces not satisfied.
	var _ ClusterDistancer = CentroidGeoDistance{}
}
