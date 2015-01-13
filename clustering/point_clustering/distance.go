package point_clustering

// A ClusterDistancer defines the how to compute the distance between point clusters.
type ClusterDistancer interface {
	ClusterDistance(pc1, pc2 *Cluster) float64
}

// ClusterDistancerFactory builds ClusterDistancers.
// Since these distancers can also hold state, a different one must be used
// for each clustering operations. That new distancer is built using this interface.
type ClusterDistancerFactory interface {
	// ClusterDistancer clusterCount and pointCount parameters can help with
	// creating the caching object, if necessary.
	ClusterDistancer(clusterCount, pointCount int) ClusterDistancer
}

// CentroidDistance implements the ClusterDistancer interface where the
// distance is just the euclidean distance between the cluster centroids.
type CentroidDistance struct{}

// ClusterDistance computes the distance between the cluster centroids.
func (cd CentroidDistance) ClusterDistance(c1, c2 *Cluster) float64 {
	return c1.Centroid.DistanceFrom(c2.Centroid)
}

// ClusterDistancer returns itself, since it is also a ClusterDistancer.
func (cd CentroidDistance) ClusterDistancer(clusterCount, pointCount int) ClusterDistancer {
	return cd
}

// CentroidSquaredDistance implements the ClusterDistancer interface where the
// distance is just the squared euclidean distance between the cluster centroids.
// This distancer is recommended over CentroidDistance, just square the threshold.
type CentroidSquaredDistance struct{}

// ClusterDistance computes the squared euclidean distance between the cluster centroids.
func (csd CentroidSquaredDistance) ClusterDistance(c1, c2 *Cluster) float64 {
	return c1.Centroid.SquaredDistanceFrom(c2.Centroid)
}

// ClusterDistancer returns itself, since it is also a ClusterDistancer.
func (csd CentroidSquaredDistance) ClusterDistancer(clusterCount, pointCount int) ClusterDistancer {
	return csd
}

// CentroidGeoDistance implements the ClusterDistancer interface where the
// distance is just the geo distance between the Group centroids.
// If possible, it is recommended to project the lat/lng points into a
// euclidean space and use CentroidSquaredDistance.
type CentroidGeoDistance struct{}

// ClusterDistance computes the geo distance between the cluster centroids.
func (cgd CentroidGeoDistance) ClusterDistance(c1, c2 *Cluster) float64 {
	return c1.Centroid.GeoDistanceFrom(c2.Centroid)
}

// ClusterDistancer returns itself, since it is also a ClusterDistancer.
func (cgd CentroidGeoDistance) ClusterDistancer(clusterCount, pointCount int) ClusterDistancer {
	return cgd
}
