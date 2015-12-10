package clustering

import "github.com/paulmach/go.geo"

// A Cluster is a cluster of pointers plus their centroid.
// It defines a center/centroid for easy centroid distance computation.
type Cluster struct {
	Centroid *geo.Point
	Pointers []geo.Pointer
}

// NewCluster creates the point cluster and finds the center of the given pointers.
func NewCluster(pointers ...geo.Pointer) *Cluster {
	var (
		sumX, sumY float64
		count      int
	)

	c := &Cluster{
		Pointers: pointers,
	}

	if len(pointers) == 0 {
		c.Centroid = geo.NewPoint(0, 0)
		return c
	}

	if len(pointers) == 1 {
		c.Centroid = pointers[0].Point().Clone()
		return c
	}

	// find the center/centroid of multiple points
	for _, pointer := range c.Pointers {
		cp := pointer.Point()

		sumX += cp.X()
		sumY += cp.Y()
		count++
	}
	c.Centroid = geo.NewPoint(sumX/float64(count), sumY/float64(count))

	return c
}

// NewClusterWithCentroid creates a point cluster stub from the given centroid
// and optional pointers.
func NewClusterWithCentroid(centroid *geo.Point, pointers ...geo.Pointer) *Cluster {
	return &Cluster{
		Centroid: centroid.Clone(),
		Pointers: pointers,
	}
}

func (c *Cluster) merge(c2 *Cluster) {

	percent := 1 - float64(len(c.Pointers))/float64(len(c.Pointers)+len(c2.Pointers))

	c.Centroid.SetX(c.Centroid[0] + percent*(c2.Centroid[0]-c.Centroid[0]))
	c.Centroid.SetY(c.Centroid[1] + percent*(c2.Centroid[1]-c.Centroid[1]))
	c.Pointers = append(c.Pointers, c2.Pointers...)
}
