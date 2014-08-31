package reducers

import (
	"github.com/paulmach/go.geo"
)

type distanceFunc func(*geo.Point, *geo.Point) float64

// A RadialReducer wraps the Radial function
// to fulfill the geo.Reducer and geo.GeoReducer interfaces.
type RadialReducer struct {
	Threshold float64 // euclidean distance
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

// GeoReduce runs the RadialGeo reduction. The path should be in lng/lat (EPSG:4326).
// The threshold is expected to be in meters.
func (r RadialReducer) GeoReduce(path *geo.Path) *geo.Path {
	return RadialGeo(path, r.Threshold)
}

// A RadialGeoReducer wraps the RadialGeo function
// to fulfill the geo.Reducer and geo.GeoReducer interfaces.
type RadialGeoReducer struct {
	Threshold float64 // meters
}

// NewRadialGeoReducer creates a new RadialGeoReducer.
// This reducer should be used with EPSG:4326 (lng/lat) paths.
func NewRadialGeoReducer(meters float64) *RadialGeoReducer {
	return &RadialGeoReducer{
		Threshold: meters,
	}
}

// Reduce runs the RadialGeo reduction using the threshold of the RadialGeoReducer.
// The threshold is expected to be in meters.
func (r RadialGeoReducer) Reduce(path *geo.Path) *geo.Path {
	return RadialGeo(path, r.Threshold)
}

// GeoReduce runs the RadialGeo reduction. The path should be in lng/lat (EPSG:4326).
// The threshold is expected to be in meters.
func (r RadialGeoReducer) GeoReduce(path *geo.Path) *geo.Path {
	return RadialGeo(path, r.Threshold)
}

// Radial peforms a radial distance polyline simplification using a standard euclidean distance.
// Returns a new path and DOES NOT modify the original.
func Radial(path *geo.Path, meters float64) *geo.Path {
	p, _ := radialCore(path, meters*meters, squaredDistance, false)
	return p
}

// RadialIndexMap is similar to Radial but returns an array that maps
// each new path index to its original path index.
// Returns a new path and DOES NOT modify the original.
func RadialIndexMap(path *geo.Path, meters float64) (*geo.Path, []int) {
	return radialCore(path, meters*meters, squaredDistance, true)
}

// RadialGeo peforms a radial distance polyline simplification using the GeoDistance.
// ie. the path points must be lng/lat points otherwise the behavior of this function is undefined.
// Returns a new path and DOES NOT modify the original.
func RadialGeo(path *geo.Path, meters float64) *geo.Path {
	p, _ := radialCore(path, meters, geoDistance, false)
	return p
}

// RadialGeoIndexMap is similar to RadialGeo but returns an array that maps
// each new path index to its original path index.
// Returns a new path and DOES NOT modify the original.
func RadialGeoIndexMap(path *geo.Path, meters float64) (*geo.Path, []int) {
	return radialCore(path, meters, geoDistance, true)
}

func radialCore(
	path *geo.Path,
	meters float64,
	dist distanceFunc,
	needIndexMap bool,
) (*geo.Path, []int) {

	// initial sanity checks
	if path.Length() == 0 {
		return path.Clone(), []int{}
	}

	if path.Length() == 1 {
		return path.Clone(), []int{0}
	}

	if path.Length() == 2 {
		return path.Clone(), []int{0, 1}
	}

	var newPoints []geo.Point
	var indexMap []int

	points := path.Points()
	newPoints = append(newPoints, points[0])

	if needIndexMap {
		indexMap = append(indexMap, 0)
	}

	// split it up this way because I think it's faster
	// TODO: test this assumption
	currentIndex := 0
	if needIndexMap {
		for i := 1; i < len(points); i++ {
			if dist(&points[currentIndex], &points[i]) > meters {
				currentIndex = i
				indexMap = append(indexMap, currentIndex)
				newPoints = append(newPoints, points[i])
			}
		}
	} else {
		for i := 1; i < len(points); i++ {
			if dist(&points[currentIndex], &points[i]) > meters {
				currentIndex = i
				newPoints = append(newPoints, points[i])
			}
		}
	}

	if currentIndex != len(points)-1 {
		newPoints = append(newPoints, points[len(points)-1])
		if needIndexMap {
			indexMap = append(indexMap, len(points)-1)
		}
	}

	p := &geo.Path{}
	return p.SetPoints(newPoints), indexMap
}

func squaredDistance(p1, p2 *geo.Point) float64 {
	return p1.SquaredDistanceFrom(p2)
}

func geoDistance(p1, p2 *geo.Point) float64 {
	return p1.GeoDistanceFrom(p2)
}
