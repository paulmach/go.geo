package path_clustering

import "testing"

func TestSingleLinkageDistancerFactory(t *testing.T) {
	// will not compile if interfaces not satisfied.
	var _ ClusterDistancerFactory = NewSingleLinkageDistancerFactory(nil)
}

func TestCompleteLinkageDistancerFactory(t *testing.T) {
	// will not compile if interfaces not satisfied.
	var _ ClusterDistancerFactory = NewCompleteLinkageDistancerFactory(nil)
}

func TestSingleLinkageDistancer(t *testing.T) {
	// will not compile if interfaces not satisfied.
	var _ ClusterDistancer = NewSingleLinkageDistancer(nil, 0, 0)
}

func TestCompleteLinkageDistancer(t *testing.T) {
	// will not compile if interfaces not satisfied.
	var _ ClusterDistancer = NewCompleteLinkageDistancer(nil, 0, 0)
}
