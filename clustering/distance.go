package clustering

// A ClusterDistancer defines the how to compute the distance between point clusters.
type ClusterDistancer interface {
	ClusterDistance(c1, c2 *Cluster) float64
}

// CentroidDistance implements the ClusterDistancer interface where the
// distance is just the euclidean distance between the cluster centroids.
type CentroidDistance struct{}

// ClusterDistance computes the distance between the cluster centroids.
func (cd CentroidDistance) ClusterDistance(c1, c2 *Cluster) float64 {
	return c1.Centroid.DistanceFrom(c2.Centroid)
}

// CentroidSquaredDistance implements the ClusterDistancer interface where the
// distance is just the squared euclidean distance between the cluster centroids.
// This distancer is recommended over CentroidDistance, just square the threshold.
type CentroidSquaredDistance struct{}

// ClusterDistance computes the squared euclidean distance between the cluster centroids.
func (csd CentroidSquaredDistance) ClusterDistance(c1, c2 *Cluster) float64 {
	// save the function call, is this faster?
	d0 := (c1.Centroid[0] - c2.Centroid[0])
	d1 := (c1.Centroid[1] - c2.Centroid[1])
	return d0*d0 + d1*d1
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
