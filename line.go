package geo

import (
	"math"
)

// represents the shortest path between A and B
type Line struct {
	A, B Point
}

// Transform applies a given projection or inverse projection to the current line.
// Modifies the line.
func (l *Line) Transform(projection func(*Point) *Point) *Line {
	projection(&l.A)
	projection(&l.B)

	return l
}

// DistanceFrom does NOT use geodesic geometry. It finds the distance from
// the line using standard euclidean geometry, using the units the points are in.
func (l *Line) DistanceFrom(point *Point) float64 {

	if l.A.Equals(&l.B) {
		// line is of length 0
		return l.A.DistanceFrom(point)
	} else {
		u := ((point.Y()-l.A.Y())*(l.B.Y()-l.A.Y()) + (point.X()-l.A.X())*(l.B.X()-l.A.X())) / (math.Pow(l.B.Y()-l.A.Y(), 2) + math.Pow(l.B.X()-l.A.X(), 2))

		if u <= 0 {
			return l.A.DistanceFrom(point)
		} else if u >= 1 {
			return l.B.DistanceFrom(point)
		} else {
			return l.Interpolate(u).DistanceFrom(point)
		}
	}
}

// Distance computes the distance of the line, ie. its length, in euclidian space.
func (l *Line) Distance() float64 {
	return l.A.DistanceFrom(&l.B)
}

// GeoDistance the distance of the line, ie. its length, using spherical geometry.
func (l *Line) GeoDistance(haversine ...bool) float64 {
	return l.A.GeoDistanceFrom(&l.B, yesHaversine(haversine))
}

// Interpolate performs a simple linear interpolation, from A to B
func (l *Line) Interpolate(percent float64) *Point {
	p := &Point{}
	p.SetX(l.A.X() + percent*(l.B.X()-l.A.X()))
	p.SetY(l.A.Y() + percent*(l.B.Y()-l.A.Y()))

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

	dLng := deg2rad(l.B.Lng() - l.A.Lng())

	aLatRad := deg2rad(l.A.Lat())
	bLatRad := deg2rad(l.B.Lat())

	x := math.Cos(bLatRad) * math.Cos(dLng)
	y := math.Cos(bLatRad) * math.Sin(dLng)

	p.SetLat(math.Atan2(math.Sin(aLatRad)+math.Sin(bLatRad), math.Sqrt((math.Cos(aLatRad)+x)*(math.Cos(aLatRad)+x)+y*y)))
	p.SetLng(deg2rad(l.A.Lng()) + math.Atan2(y, math.Cos(aLatRad)+x))

	// convert back to degrees
	p.SetLat(rad2deg(p.Lat()))
	p.SetLng(rad2deg(p.Lng()))

	return p
}

// Bounds returns bound around the line. Simply uses rectangular coordinates.
func (l *Line) Bounds() *Bound {
	return NewBound(math.Max(l.A.X(), l.B.X()), math.Min(l.A.X(), l.B.X()),
		math.Max(l.A.Y(), l.B.Y()), math.Min(l.A.Y(), l.B.Y()))
}

// Reverse swapps the start and end of the line
func (l *Line) Reverse() *Line {
	l.A, l.B = l.B, l.A
	return l
}

func (l *Line) Clone() *Line {
	return &Line{*l.A.Clone(), *l.B.Clone()}
}
