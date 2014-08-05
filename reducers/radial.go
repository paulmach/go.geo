package reducers

import (
	"github.com/paulmach/go.geo"
)

type distanceFunc func(*geo.Point, *geo.Point) float64

// A RadialReducer wraps the Radial function
// to fulfil the geo.Reducer interface.
type RadialReducer struct {
	Threshold float64 // meters
}

// NewRadialReducer creates a new RadialReducer.
func NewRadialReducer(threshold float64) *RadialReducer {
	return &RadialReducer{
		Threshold: threshold,
	}
}

// Reduce runs the Radial reduction using the threshold of the RadialReducer.
func (r RadialReducer) Reduce(path *geo.Path) *geo.Path {
	return Radial(path, r.Threshold)
}

// A RadialGeoReducer wraps the RadialGeo function
// to fulfil the geo.Reducer interface.
type RadialGeoReducer struct {
	Threshold float64 // meters
}

// NewRadialGeoReducer creates a new RadialGeoReducer.
func NewRadialGeoReducer(meters float64) *RadialGeoReducer {
	return &RadialGeoReducer{
		Threshold: meters,
	}
}

// Reduce runs the RadialGeo reduction using the threshold of the RadialGeoReducer.
func (r RadialGeoReducer) Reduce(path *geo.Path) *geo.Path {
	return RadialGeo(path, r.Threshold)
}

// Radial peforms a radial distance polyline simplification using a standard euclidean distance.
// Returns a new path and DOES NOT modify the original.
func Radial(path *geo.Path, meters float64) *geo.Path {
	return radialCore(path, meters, distance)
}

// RadialGeo peforms a radial distance polyline simplification using the GeoDistance.
// ie. the path points must be lng/lat points otherwise the behavior of this function is undefined.
// Returns a new path and DOES NOT modify the original.
func RadialGeo(path *geo.Path, meters float64) *geo.Path {
	return radialCore(path, meters, geoDistance)
}

func radialCore(path *geo.Path, meters float64, dist distanceFunc) *geo.Path {
	if path.Length() == 0 {
		return path
	}

	mask := make([]byte, path.Length())
	mask[0] = 1
	mask[path.Length()-1] = 1

	points := path.Points()
	newPoints := make([]geo.Point, 1, len(points)/2+1)
	newPoints[0] = *points[0].Clone()

	currentIndex := 0
	for i := 1; i < len(points); i++ {
		if dist(&points[currentIndex], &points[i]) > meters {
			currentIndex = i
			newPoints = append(newPoints, points[i])
		}
	}

	if currentIndex != len(points)-1 {
		newPoints = append(newPoints, points[len(points)-1])
	}

	p := &geo.Path{}
	return p.SetPoints(newPoints)
}

func distance(p1, p2 *geo.Point) float64 {
	return p1.DistanceFrom(p2)
}

func geoDistance(p1, p2 *geo.Point) float64 {
	return p1.GeoDistanceFrom(p2)
}
