package reducers

import (
	"github.com/paulmach/go.geo"
)

type distanceFunc func(*geo.Point, *geo.Point) float64

func distance(p1, p2 *geo.Point) float64 {
	return p1.DistanceFrom(p2)
}

func geoDistance(p1, p2 *geo.Point) float64 {
	return p1.GeoDistanceFrom(p2)
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
