package geo

import (
	"fmt"
	"math"
)

// A simple X/Y or Lng/Lat 2d point
type Point [2]float64 // [X, Y] or [Lng, Lat]
var InfinityPoint = &Point{math.Inf(1), math.Inf(1)}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

// Transform applies a given projection or inverse projection to the current point.
func (p *Point) Transform(transformer func(*Point) *Point) *Point {
	transformer(p)
	return p
}

func (p *Point) DistanceFrom(point *Point) float64 {
	d0 := (point[0] - p[0])
	d1 := (point[1] - p[1])
	return math.Sqrt(d0*d0 + d1*d1)
}

// GeoDistanceFrom returns the geodesic distance in meters
func (p *Point) GeoDistanceFrom(point *Point, haversine ...bool) float64 {
	dLat := deg2rad(point.Lat() - p.Lat())
	dLng := deg2rad(point.Lng() - p.Lng())

	if yesHaversine(haversine) {
		// yes trig functions
		dLat2Sin := math.Sin(dLat / 2)
		dLng2Sin := math.Sin(dLng / 2)
		a := dLat2Sin*dLat2Sin + math.Cos(deg2rad(p.Lat()))*math.Cos(deg2rad(point.Lat()))*dLng2Sin*dLng2Sin

		return 2.0 * DEFAULT_Radius * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	}

	// fast way using pythagorean theorem on an equirectangular projection
	x := dLng * math.Cos(deg2rad((p.Lat()+point.Lat())/2.0))
	return math.Sqrt(dLat*dLat+x*x) * DEFAULT_Radius
}

// WARNING: untested
func (p *Point) BearingTo(point *Point) float64 {
	dLng := deg2rad(point.Lng() - p.Lng())

	pLatRad := deg2rad(p.Lat())
	pointLatRad := deg2rad(point.Lat())

	y := math.Sin(dLng) * math.Cos(point.Lat())
	x := math.Cos(pLatRad)*math.Sin(pointLatRad) - math.Sin(pLatRad)*math.Cos(pointLatRad)*math.Cos(dLng)

	return rad2deg(math.Atan2(y, x))
}

// Add a point to the given point
func (p *Point) Add(point *Point) *Point {
	p[0] += point[0]
	p[1] += point[1]

	return p
}

// Subtract a point from the given point
func (p *Point) Subtract(point *Point) *Point {
	p[0] -= point[0]
	p[1] -= point[1]

	return p
}

// Normalize treats the point as a vector and
// scales it such that its distance from [0,0] is 1
func (p *Point) Normalize() *Point {
	dist := p.DistanceFrom(&Point{})

	p[0] /= dist
	p[1] /= dist

	return p
}

// Scale each component of the point
func (p *Point) Scale(factor float64) *Point {
	p[0] *= factor
	p[1] *= factor

	return p
}

func (p *Point) ToArray() [2]float64 {
	return [2]float64(*p)
}

func (p *Point) Clone() *Point {
	newP := &Point{}
	copy(newP[:], p[:])
	return newP
}

func (p *Point) Equals(point *Point) bool {
	if p[0] == point[0] && p[1] == point[1] {
		return true
	}

	return false
}

func (p *Point) Lat() float64 {
	return p[1]
}

func (p *Point) SetLat(lat float64) *Point {
	p[1] = lat
	return p
}

func (p *Point) Lng() float64 {
	return p[0]
}

func (p *Point) SetLng(lng float64) *Point {
	p[0] = lng
	return p
}

func (p *Point) X() float64 {
	return p[0]
}

func (p *Point) SetX(x float64) *Point {
	p[0] = x
	return p
}

func (p *Point) Y() float64 {
	return p[1]
}

func (p *Point) SetY(y float64) *Point {
	p[1] = y
	return p
}

func (p *Point) String() string {
	return fmt.Sprintf("[%f, %f]", p[0], p[1])
}
