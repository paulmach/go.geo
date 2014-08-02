package geo

import (
	"fmt"
	"math"
)

// A Point is a simple X/Y or Lng/Lat 2d point. [X, Y] or [Lng, Lat]
type Point [2]float64

// InfinityPoint is the point at [inf, inf].
// Currently returned for the intersection of two collinear overlapping lines.
var InfinityPoint = &Point{math.Inf(1), math.Inf(1)}

// NewPoint creates a new point
func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

// Transform applies a given projection or inverse projection to the current point.
func (p *Point) Transform(projector Projector) *Point {
	projector(p)
	return p
}

// DistanceFrom returns the Euclidean distance between the points.
func (p *Point) DistanceFrom(point *Point) float64 {
	d0 := (point[0] - p[0])
	d1 := (point[1] - p[1])
	return math.Sqrt(d0*d0 + d1*d1)
}

// GeoDistanceFrom returns the geodesic distance in meters.
func (p *Point) GeoDistanceFrom(point *Point, haversine ...bool) float64 {
	dLat := deg2rad(point.Lat() - p.Lat())
	dLng := deg2rad(point.Lng() - p.Lng())

	if yesHaversine(haversine) {
		// yes trig functions
		dLat2Sin := math.Sin(dLat / 2)
		dLng2Sin := math.Sin(dLng / 2)
		a := dLat2Sin*dLat2Sin + math.Cos(deg2rad(p.Lat()))*math.Cos(deg2rad(point.Lat()))*dLng2Sin*dLng2Sin

		return 2.0 * EarthRadius * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	}

	// fast way using pythagorean theorem on an equirectangular projection
	x := dLng * math.Cos(deg2rad((p.Lat()+point.Lat())/2.0))
	return math.Sqrt(dLat*dLat+x*x) * EarthRadius
}

// BearingTo computes the direction one must start traveling on earth
// to be heading to the given point. WARNING: untested
func (p *Point) BearingTo(point *Point) float64 {
	dLng := deg2rad(point.Lng() - p.Lng())

	pLatRad := deg2rad(p.Lat())
	pointLatRad := deg2rad(point.Lat())

	y := math.Sin(dLng) * math.Cos(point.Lat())
	x := math.Cos(pLatRad)*math.Sin(pointLatRad) - math.Sin(pLatRad)*math.Cos(pointLatRad)*math.Cos(dLng)

	return rad2deg(math.Atan2(y, x))
}

// Add a point to the given point.
func (p *Point) Add(point *Point) *Point {
	p[0] += point[0]
	p[1] += point[1]

	return p
}

// Subtract a point from the given point.
func (p *Point) Subtract(point *Point) *Point {
	p[0] -= point[0]
	p[1] -= point[1]

	return p
}

// Normalize treats the point as a vector and
// scales it such that its distance from [0,0] is 1.
func (p *Point) Normalize() *Point {
	dist := p.DistanceFrom(&Point{})

	if dist == 0 {
		p[0] = 0
		p[1] = 0

		return p
	}

	p[0] /= dist
	p[1] /= dist

	return p
}

// Scale each component of the point.
func (p *Point) Scale(factor float64) *Point {
	p[0] *= factor
	p[1] *= factor

	return p
}

// Dot is just x1*x2 + y1*y2
func (p *Point) Dot(v *Point) float64 {
	return p[0]*v[0] + p[1]*v[1]
}

// ToArray casts the data to a [2]float64.
func (p *Point) ToArray() [2]float64 {
	return [2]float64(*p)
}

// Clone creates a duplicate of the point.
func (p *Point) Clone() *Point {
	newP := &Point{}
	copy(newP[:], p[:])
	return newP
}

// Equals checks if the point represents the same point or vector.
func (p *Point) Equals(point *Point) bool {
	if p[0] == point[0] && p[1] == point[1] {
		return true
	}

	return false
}

// Lat returns the latitude/vertical component of the point.
func (p *Point) Lat() float64 {
	return p[1]
}

// SetLat sets the latitude/vertical component of the point.
func (p *Point) SetLat(lat float64) *Point {
	p[1] = lat
	return p
}

// Lng returns the longitude/horizontal component of the point.
func (p *Point) Lng() float64 {
	return p[0]
}

// SetLng sets the longitude/horizontal component of the point.
func (p *Point) SetLng(lng float64) *Point {
	p[0] = lng
	return p
}

// X returns the x/horizontal component of the point.
func (p *Point) X() float64 {
	return p[0]
}

// SetX sets the x/horizontal component of the point.
func (p *Point) SetX(x float64) *Point {
	p[0] = x
	return p
}

// Y returns the y/vertical component of the point.
func (p *Point) Y() float64 {
	return p[1]
}

// SetY sets the y/vertical component of the point.
func (p *Point) SetY(y float64) *Point {
	p[1] = y
	return p
}

// String returns a string representation of the point.
func (p *Point) String() string {
	return fmt.Sprintf("[%f, %f]", p[0], p[1])
}
