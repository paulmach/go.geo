package geo

import (
	"math"
)

// represents the shortest path between A and B
type Line struct {
	a, b Point
}

func NewLine(a, b *Point) *Line {
	return &Line{*a.Clone(), *b.Clone()}
}

// Transform applies a given projection or inverse projection to the current line.
// Modifies the line.
func (l *Line) Transform(projection func(*Point) *Point) *Line {
	projection(&l.a)
	projection(&l.b)

	return l
}

// DistanceFrom does NOT use geodesic geometry. It finds the distance from
// the line using standard euclidean geometry, using the units the points are in.
func (l *Line) DistanceFrom(point *Point) float64 {

	if l.a.Equals(&l.b) {
		// line is of length 0
		return l.a.DistanceFrom(point)
	} else {
		u := ((point.Y()-l.a.Y())*(l.b.Y()-l.a.Y()) + (point.X()-l.a.X())*(l.b.X()-l.a.X())) / (math.Pow(l.b.Y()-l.a.Y(), 2) + math.Pow(l.b.X()-l.a.X(), 2))

		if u <= 0 {
			return l.a.DistanceFrom(point)
		} else if u >= 1 {
			return l.b.DistanceFrom(point)
		} else {
			return l.Interpolate(u).DistanceFrom(point)
		}
	}
}

// Distance computes the distance of the line, ie. its length, in euclidian space.
func (l *Line) Distance() float64 {
	return l.a.DistanceFrom(&l.b)
}

// GeoDistance the distance of the line, ie. its length, using spherical geometry.
func (l *Line) GeoDistance(haversine ...bool) float64 {
	return l.a.GeoDistanceFrom(&l.b, yesHaversine(haversine))
}

// Interpolate performs a simple linear interpolation, from A to B
func (l *Line) Interpolate(percent float64) *Point {
	p := &Point{}
	p.SetX(l.a.X() + percent*(l.b.X()-l.a.X()))
	p.SetY(l.a.Y() + percent*(l.b.Y()-l.a.Y()))

	// simple
	return p
}

// Midpoint returns the euclidean midpoint of the line
func (l *Line) Midpoint() *Point {
	return l.Interpolate(0.5)
}

// GeoMidpoint returns the half-way point along a great circle path between the two points.
// WARNING: untested
func (l *Line) GeoMidpoint() *Point {
	p := &Point{}

	dLng := deg2rad(l.b.Lng() - l.a.Lng())

	aLatRad := deg2rad(l.a.Lat())
	bLatRad := deg2rad(l.b.Lat())

	x := math.Cos(bLatRad) * math.Cos(dLng)
	y := math.Cos(bLatRad) * math.Sin(dLng)

	p.SetLat(math.Atan2(math.Sin(aLatRad)+math.Sin(bLatRad), math.Sqrt((math.Cos(aLatRad)+x)*(math.Cos(aLatRad)+x)+y*y)))
	p.SetLng(deg2rad(l.a.Lng()) + math.Atan2(y, math.Cos(aLatRad)+x))

	// convert back to degrees
	p.SetLat(rad2deg(p.Lat()))
	p.SetLng(rad2deg(p.Lng()))

	return p
}

// Bounds returns bound around the line. Simply uses rectangular coordinates.
func (l *Line) Bounds() *Bound {
	return NewBound(math.Max(l.a.X(), l.b.X()), math.Min(l.a.X(), l.b.X()),
		math.Max(l.a.Y(), l.b.Y()), math.Min(l.a.Y(), l.b.Y()))
}

// Reverse swapps the start and end of the line
func (l *Line) Reverse() *Line {
	l.a, l.b = l.b, l.a
	return l
}

// direction is irrelevant, ie. true if one is the reverse of the other
func (l *Line) Equals(line *Line) bool {
	return (l.a.Equals(&line.a) && l.b.Equals(&line.b)) || (l.a.Equals(&line.b) && l.b.Equals(&line.a))
}

func (l *Line) Clone() *Line {
	return &Line{*l.a.Clone(), *l.b.Clone()}
}

func (l *Line) A() *Point {
	return &l.a
}

func (l *Line) B() *Point {
	return &l.b
}
