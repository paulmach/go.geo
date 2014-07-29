package geo

import (
	"math"
)

// Line represents the shortest path between A and B.
type Line struct {
	a, b Point
}

// NewLine creates a new line by cloning the provided points.
func NewLine(a, b *Point) *Line {
	return &Line{*a.Clone(), *b.Clone()}
}

// Transform applies a given projection or inverse projection to the current line.
// Modifies the line.
func (l *Line) Transform(projector Projector) *Line {
	projector(&l.a)
	projector(&l.b)

	return l
}

// DistanceFrom does NOT use geodesic geometry. It finds the distance from
// the line using standard Euclidean geometry, using the units the points are in.
func (l *Line) DistanceFrom(point *Point) float64 {

	if l.a.Equals(&l.b) {
		// line is of length 0
		return l.a.DistanceFrom(point)
	}

	u := l.Project(point)
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

// Direction computes the direction the line is pointing from A() to B().
// The units are radians from the positive x-axis.
// Range same as math.Atan2, [-Pi, Pi]
func (l *Line) Direction() float64 {
	return math.Atan2(l.b[1]-l.a[1], l.b[0]-l.a[0])
}

// Project returns the normalized distance of the point on the line nearest the given point.
// Returned values maybe the outside of [0,1]. This function is the opposite of Interpolate.
func (l *Line) Project(point *Point) float64 {
	if point.Equals(&l.a) {
		return 0.0
	}

	if point.Equals(&l.b) {
		return 1.0
	}

	dx := l.b[0] - l.a[0]
	dy := l.b[1] - l.a[1]
	return ((point[0]-l.a[0])*dx + (point[1]-l.a[1])*dy) / (dx*dx + dy*dy)
}

// Measure returns the distance along the line to the point nearest the given point.
// Treats the line as a line segment such that is the nearest point is an endpoint of the line,
// the function will return 0 or 1 as appropriate.
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
// This function is the opposite of Project.
func (l *Line) Interpolate(percent float64) *Point {
	p := &Point{}
	p.SetX(l.a[0] + percent*(l.b[0]-l.a[0]))
	p.SetY(l.a[1] + percent*(l.b[1]-l.a[1]))

	// simple
	return p
}

// Side returns 1 if the point is on the right side, -1 if on the left side, and 0 if collinear.
func (l *Line) Side(p *Point) int {
	val := (l.b[0]-l.a[0])*(p[1]-l.b[1]) - (l.b[1]-l.a[1])*(p[0]-l.b[0])

	if val < 0 {
		return 1 // right
	} else if val > 0 {
		return -1 // left
	}

	return 0 // collinear
}

// Intersection finds the intersection of the two lines or nil,
// if the lines are collinear will return NewPoint(math.Inf(1), math.Inf(1)) == InfinityPoint
func (l *Line) Intersection(line *Line) *Point {
	den := (line.b[1]-line.a[1])*(l.b[0]-l.a[0]) - (line.b[0]-line.a[0])*(l.b[1]-l.a[1])
	U1 := (line.b[0]-line.a[0])*(l.a[1]-line.a[1]) - (line.b[1]-line.a[1])*(l.a[0]-line.a[0])
	U2 := (l.b[0]-l.a[0])*(l.a[1]-line.a[1]) - (l.b[1]-l.a[1])*(l.a[0]-line.a[0])

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

	return l.Interpolate(U1 / den)
}

// Intersects will return true if the lines are collinear AND intersect.
// Based on: http://www.geeksforgeeks.org/check-if-two-given-line-segments-intersect/
func (l *Line) Intersects(line *Line) bool {
	s1 := l.Side(&line.a)
	s2 := l.Side(&line.b)
	s3 := line.Side(&l.a)
	s4 := line.Side(&l.b)

	if s1 != s2 && s3 != s4 {
		return true
	}

	// Special Cases
	// l1 and l2.a collinear, check if l2.a is on l1
	lBound := l.Bounds()
	if s1 == 0 && lBound.Contains(&line.a) {
		return true
	}

	// l1 and l2.b collinear, check if l2.b is on l1
	if s2 == 0 && lBound.Contains(&line.b) {
		return true
	}

	// TODO: are these next two tests redudant give the test above.
	// Thinking yes if there is round off magic.

	// l2 and l1.a collinear, check if l1.a is on l2
	lineBound := line.Bounds()
	if s3 == 0 && lineBound.Contains(&l.a) {
		return true
	}

	// l2 and l1.b collinear, check if l1.b is on l2
	if s4 == 0 && lineBound.Contains(&l.b) {
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
	return NewBound(math.Max(l.a[0], l.b[0]), math.Min(l.a[0], l.b[0]),
		math.Max(l.a[1], l.b[1]), math.Min(l.a[1], l.b[1]))
}

// Reverse swaps the start and end of the line.
func (l *Line) Reverse() *Line {
	l.a, l.b = l.b, l.a
	return l
}

// Equals returns the line equality and is irrespective of direction,
// i.e. true if one is the reverse of the other.
func (l *Line) Equals(line *Line) bool {
	return (l.a.Equals(&line.a) && l.b.Equals(&line.b)) || (l.a.Equals(&line.b) && l.b.Equals(&line.a))
}

// Clone returns a deep copy of the line.
func (l *Line) Clone() *Line {
	return &Line{*l.a.Clone(), *l.b.Clone()}
}

// A returns a pointer to the first point in the line.
func (l *Line) A() *Point {
	return &l.a
}

// B returns a pointer to the second point in the line.
func (l *Line) B() *Point {
	return &l.b
}
