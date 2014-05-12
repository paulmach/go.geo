package geo

import (
	"math"
)

// Line represents the shortest path between A and B.
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
// the line using standard Euclidean geometry, using the units the points are in.
func (l *Line) DistanceFrom(point *Point) float64 {

	if l.a.Equals(&l.b) {
		// line is of length 0
		return l.a.DistanceFrom(point)
	}

	u := ((point.Y()-l.a.Y())*(l.b.Y()-l.a.Y()) + (point.X()-l.a.X())*(l.b.X()-l.a.X())) / (math.Pow(l.b.Y()-l.a.Y(), 2) + math.Pow(l.b.X()-l.a.X(), 2))
	if u <= 0 {
		return l.a.DistanceFrom(point)
	} else if u >= 1 {
		return l.b.DistanceFrom(point)
	}

	return l.Interpolate(u).DistanceFrom(point)
}

// Distance computes the distance of the line, ie. its length, in Euclidian space.
func (l *Line) Distance() float64 {
	return l.a.DistanceFrom(&l.b)
}

// GeoDistance computes the distance of the line, ie. its length, using spherical geometry.
func (l *Line) GeoDistance(haversine ...bool) float64 {
	return l.a.GeoDistanceFrom(&l.b, yesHaversine(haversine))
}

// Project computes the factor to multiply the line by to be nearest the given point.
func (l *Line) Project(point *Point) float64 {
	if point.Equals(l.A()) {
		return 0.0
	}
	if point.Equals(l.B()) {
		return 1.0
	}
	dx := l.B().X() - l.A().X()
	dy := l.B().Y() - l.A().Y()
	sq := dx*dx + dy*dy
	p := ((point.X()-l.A().X())*dx + (point.Y()-l.A().Y())*dy) / sq
	return p
}

// Measure computes the distance along this line to the point nearest the given point.
func (l *Line) Measure(point *Point) float64 {
	projFactor := l.Project(point)
	if projFactor <= 0.0 {
		return 0.0
	}
	if projFactor <= 1.0 {
		return projFactor * l.Distance()
	}
	// projFactor is > 1
	return l.Distance()
}

// Interpolate performs a simple linear interpolation, from A to B.
func (l *Line) Interpolate(percent float64) *Point {
	p := &Point{}
	p.SetX(l.a.X() + percent*(l.b.X()-l.a.X()))
	p.SetY(l.a.Y() + percent*(l.b.Y()-l.a.Y()))

	// simple
	return p
}

// Side returns 1 if the point is on the right side, -1 if on the left side, and 0 if collinear.
func (l *Line) Side(p *Point) int {
	val := (l.b.X()-l.a.X())*(p.Y()-l.b.Y()) - (l.b.Y()-l.a.Y())*(p.X()-l.b.X())

	if val < 0 {
		return 1 // right
	} else if val > 0 {
		return -1 // left
	}

	return 0 // collinear
}

// Intersection finds the intersection of the two lines or nil,
// if the lines are collinear will return NewPoint(math.Inf(1), math.Inf(1)) == InfinityPoint
func (l1 *Line) Intersection(l2 *Line) *Point {
	den := (l2.b.Y()-l2.a.Y())*(l1.b.X()-l1.a.X()) - (l2.b.X()-l2.a.X())*(l1.b.Y()-l1.a.Y())
	U1 := (l2.b.X()-l2.a.X())*(l1.a.Y()-l2.a.Y()) - (l2.b.Y()-l2.a.Y())*(l1.a.X()-l2.a.X())
	U2 := (l1.b.X()-l1.a.X())*(l1.a.Y()-l2.a.Y()) - (l1.b.Y()-l1.a.Y())*(l1.a.X()-l2.a.X())

	if den == 0 {
		// collinear, all bets are off
		if U1 == 0 && U2 == 0 {
			return InfinityPoint
		}

		return nil
	}

	if U1/den < 0 || U1/den > 1 || U2/den < 0 || U2/den > 1 {
		return nil
	}

	return l1.Interpolate(U1 / den)
}

// Intersects will return true if the lines are collinear AND intersect.
// Based on: http://www.geeksforgeeks.org/check-if-two-given-line-segments-intersect/
func (l1 *Line) Intersects(l2 *Line) bool {
	s1 := l1.Side(l2.A())
	s2 := l1.Side(l2.B())
	s3 := l2.Side(l1.A())
	s4 := l2.Side(l1.B())

	if s1 != s2 && s3 != s4 {
		return true
	}

	// Special Cases
	// l1 and l2.a collinear, check if l2.a is on l1
	if s1 == 0 && l1.Bounds().Contains(l2.A()) {
		return true
	}

	// l1 and l2.b collinear, check if l2.b is on l1
	if s2 == 0 && l1.Bounds().Contains(l2.B()) {
		return true
	}

	// TODO: are these next two tests redudant give the test above

	// l2 and l1.a collinear, check if l1.a is on l2
	if s3 == 0 && l2.Bounds().Contains(l1.A()) {
		return true
	}

	// l2 and l1.b collinear, check if l1.b is on l2
	if s4 == 0 && l2.Bounds().Contains(l1.B()) {
		return true
	}

	return false
}

// Midpoint returns the Euclidean midpoint of the line.
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

// Bounds returns a bound around the line. Simply uses rectangular coordinates.
func (l *Line) Bounds() *Bound {
	return NewBound(math.Max(l.a.X(), l.b.X()), math.Min(l.a.X(), l.b.X()),
		math.Max(l.a.Y(), l.b.Y()), math.Min(l.a.Y(), l.b.Y()))
}

// Reverse swaps the start and end of the line.
func (l *Line) Reverse() *Line {
	l.a, l.b = l.b, l.a
	return l
}

// Line equality is irrespective of direction, i.e. true if one is the reverse of the other.
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
