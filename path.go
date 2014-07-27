package geo

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

// Path represents a set of points to be thought of as a polyline.
type Path struct {
	points []Point
}

// NewPath simply creates a new path.
func NewPath() *Path {
	p := &Path{}
	p.points = make([]Point, 0, 1000)

	return p
}

// SetPoints allows you to set the complete pointset yourself.
// Note that the input is an array of Points (not pointers to points).
func (p *Path) SetPoints(points []Point) *Path {
	p.points = points
	return p
}

// Points returns the raw points storred with the path
// Note the output is an array of Points (not pointers to points).
func (p *Path) Points() []Point {
	return p.points
}

// Transform applies a given projection or inverse projection to all
// the points in the path.
func (p *Path) Transform(projector Projector) *Path {
	for i := range p.points {
		projector(&p.points[i])
	}

	return p
}

// Resample converts the path into totalPoints-1 evenly spaced segments.
func (p *Path) Resample(totalPoints int) *Path {
	// degenerate case
	if len(p.points) <= 1 {
		return p
	}

	if totalPoints <= 0 {
		p.points = make([]Point, 0)
		return p
	}

	points := make([]Point, 1, totalPoints)
	points[0] = p.points[0] // start stays the same

	// location on the original line
	prevIndex := 0
	prevDistance := 0.0

	// first distance we're looking for
	step := 1
	totalDistance := p.Distance()
	currentDistance := totalDistance * float64(step) / float64(totalPoints-1)

	for {
		currentLine := NewLine(&p.points[prevIndex], &p.points[prevIndex+1])
		currentLineDistance := currentLine.Distance()
		nextDistance := prevDistance + currentLineDistance

		for currentDistance <= nextDistance {
			// need to add a point
			percent := (currentDistance - prevDistance) / currentLineDistance
			points = append(points, *currentLine.Interpolate(percent))

			// move to the next distance we want
			step++
			currentDistance = totalDistance * float64(step) / float64(totalPoints-1)
			if step == totalPoints-1 { // weird round off error on my machine
				currentDistance = totalDistance
			}
		}

		// past the current point in the original line, so move to the next one
		prevIndex++
		prevDistance = nextDistance

		if prevIndex == len(p.points)-1 {
			break
		}
	}

	// end stays the same, to handle round off errors
	if totalPoints != 1 { // for 1, we want the first point
		points[totalPoints-1] = p.points[len(p.points)-1]
	}
	p.points = points
	return p
}

// Decode is the inverse of Encode. It takes a string encoding of path
// and returns the actual path it represents. Factor defaults to 1.0e5,
// the same used by Google for polyline encoding.
func Decode(encoded string, factor ...int) *Path {
	var count, index int

	f := 1.0e5
	if len(factor) != 0 {
		f = float64(factor[0])
	}

	p := NewPath()
	tempLatLng := [2]int{0, 0}

	for index < len(encoded) {
		var result int
		var b = 0x20
		var shift uint

		for b >= 0x20 {
			b = int(encoded[index]) - 63
			index++

			result |= (b & 0x1f) << shift
			shift += 5
		}

		// sign dection
		if result&1 != 0 {
			result = ^(result >> 1)
		} else {
			result = result >> 1
		}

		if count%2 == 0 {
			result += tempLatLng[0]
			tempLatLng[0] = result
		} else {
			result += tempLatLng[1]
			tempLatLng[1] = result

			p.Push(&Point{float64(tempLatLng[1]) / f, float64(tempLatLng[0]) / f})
		}

		count++
	}

	return p
}

// Encode converts the path to a string using the Google Maps Polyline Encoding method.
// Factor defaults to 1.0e5, the same used by Google for polyline encoding.
func (p *Path) Encode(factor ...int) string {
	f := 1.0e5
	if len(factor) != 0 {
		f = float64(factor[0])
	}

	var pLat int
	var pLng int

	var result bytes.Buffer

	for _, p := range p.points {
		lat5 := int(p.Lat() * f)
		lng5 := int(p.Lng() * f)

		deltaLat := lat5 - pLat
		deltaLng := lng5 - pLng

		pLat = lat5
		pLng = lng5

		result.WriteString(encodeSignedNumber(deltaLat))
		result.WriteString(encodeSignedNumber(deltaLng))
	}

	return result.String()
}

func encodeSignedNumber(num int) string {
	shiftedNum := num << 1

	if num < 0 {
		shiftedNum = ^shiftedNum
	}

	return encodeNumber(shiftedNum)
}

func encodeNumber(num int) string {
	result := ""

	for num >= 0x20 {
		result += string((0x20 | (num & 0x1f)) + 63)
		num >>= 5
	}

	result += string(num + 63)

	return result
}

// Distance computes the total distance in the units of the points.
func (p *Path) Distance() float64 {
	sum := 0.0

	loopTo := len(p.points) - 1
	for i := 0; i < loopTo; i++ {
		sum += p.points[i].DistanceFrom(&p.points[i+1])
	}

	return sum
}

// GeoDistance computes the total distance using spherical geometry.
func (p *Path) GeoDistance(haversine ...bool) float64 {
	yesgeo := yesHaversine(haversine)
	sum := 0.0

	loopTo := len(p.points) - 1
	for i := 0; i < loopTo; i++ {
		sum += p.points[i].GeoDistanceFrom(&p.points[i+1], yesgeo)
	}

	return sum
}

// DistanceFrom computes an O(n) distance from the path. Loops over every
// subline to find the minimum distance.
func (p *Path) DistanceFrom(point *Point) float64 {
	dist := math.Inf(1)

	loopTo := len(p.points) - 1
	for i := 0; i < loopTo; i++ {
		l := &Line{p.points[i], p.points[i+1]}
		dist = math.Min(l.DistanceFrom(point), dist)
	}

	return dist
}

// Measure computes the distance along this path to the point nearest the given point.
func (p *Path) Measure(point *Point) float64 {
	minDistance := math.Inf(1)
	measure := math.Inf(-1)
	sum := 0.0
	for i := 0; i < len(p.points)-1; i++ {
		seg := &Line{p.points[i], p.points[i+1]}
		distanceToLine := seg.DistanceFrom(point)
		if distanceToLine < minDistance {
			minDistance = distanceToLine
			measure = sum + seg.Measure(point)
		}
		sum += seg.Distance()
	}
	return measure
}

// Project computes the measure along this path closest to the given point,
// normalized to the length of the path.
func (p *Path) Project(point *Point) float64 {
	return p.Measure(point) / p.Distance()
}

// Intersection calls IntersectionPath or IntersectionLine depending on the
// type of the provided geometry.
func (p *Path) Intersection(geometry interface{}) ([]*Point, [][2]int) {
	switch g := geometry.(type) {
	case Line:
		return p.IntersectionLine(&g)
	case *Line:
		return p.IntersectionLine(g)
	case Path:
		return p.IntersectionPath(&g)
	case *Path:
		return p.IntersectionPath(g)
	default:
		panic("can only determine intersection with lines and paths")
	}

	return nil, nil // unreachable
}

// IntersectionPath returns a slice of points and a slice of tuples [i, j] where i is the segment
// in the parent path and j is the segment in the given path that intersect to form the given point.
// Slices will be empty if there is no intersection.
func (p *Path) IntersectionPath(path *Path) ([]*Point, [][2]int) {
	// TODO: done some sort of line sweep here if p.Length() is big enough
	var points []*Point
	var indexes [][2]int

	for i := 0; i < len(p.points)-1; i++ {
		pLine := NewLine(&p.points[i], &p.points[i+1])

		for j := 0; j < len(path.points)-1; j++ {
			pathLine := NewLine(&path.points[j], &path.points[j+1])

			if point := pLine.Intersection(pathLine); point != nil {
				points = append(points, point)
				indexes = append(indexes, [2]int{i, j})
			}
		}
	}

	return points, indexes
}

// IntersectionLine returns a slice of points and a slice of tuples [i, 0] where i is the segment
// in path that intersects with the line at the given point.
// Slices will be empty if there is no intersection.
func (p *Path) IntersectionLine(line *Line) ([]*Point, [][2]int) {
	var points []*Point
	var indexes [][2]int

	for i := 0; i < len(p.points)-1; i++ {
		pTest := NewLine(&p.points[i], &p.points[i+1])
		if point := pTest.Intersection(line); point != nil {
			points = append(points, point)
			indexes = append(indexes, [2]int{i, 0})
		}
	}

	return points, indexes
}

// Intersects can take a line or a path to determine if there is an intersection.
func (p *Path) Intersects(geometry interface{}) bool {
	switch g := geometry.(type) {
	case Line:
		return p.IntersectsLine(&g)
	case *Line:
		return p.IntersectsLine(g)
	case Path:
		return p.IntersectsPath(&g)
	case *Path:
		return p.IntersectsPath(g)
	default:
		panic("can only determine intersection with lines and paths")
	}

	return false // unreachable
}

// IntersectsPath takes a Path and checks if it intersects with the path.
func (p *Path) IntersectsPath(path *Path) bool {
	// TODO: done some sort of line sweep here if p.Length() is big enough
	for i := 0; i < len(p.points)-1; i++ {
		pLine := NewLine(&p.points[i], &p.points[i+1])

		for j := 0; j < len(path.points)-1; j++ {
			pathLine := NewLine(&path.points[j], &path.points[j+1])

			if pLine.Intersects(pathLine) {
				return true
			}
		}
	}

	return false
}

// IntersectsLine takes a Line and checks if it intersects with the path.
func (p *Path) IntersectsLine(line *Line) bool {
	for i := 0; i < len(p.points)-1; i++ {
		pTest := NewLine(&p.points[i], &p.points[i+1])
		if pTest.Intersects(line) {
			return true
		}
	}

	return false
}

// Bound returns a bound around the path. Simply uses rectangular coordinates.
func (p *Path) Bound() *Bound {
	if len(p.points) == 0 {
		return NewBound(0, 0, 0, 0)
	}

	minX := math.Inf(1)
	minY := math.Inf(1)

	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	for _, v := range p.points {
		minX = math.Min(minX, v.X())
		minY = math.Min(minY, v.Y())

		maxX = math.Max(maxX, v.X())
		maxY = math.Max(maxY, v.Y())
	}

	return NewBound(maxX, minX, maxY, minY)
}

// SetAt updates a position at i along the path.
// Panics if index is out of range.
func (p *Path) SetAt(index int, point *Point) *Path {
	if index >= len(p.points) || index < 0 {
		panic(fmt.Sprintf("geo: set index out of range, requested: %d, length: %d", index, len(p.points)))
	}
	p.points[index] = *point
	return p
}

// GetAt returns the pointer to the Point in the page.
// This function is good for modifying values in place.
// Returns nil if index is out of range.
func (p *Path) GetAt(i int) *Point {
	if i >= len(p.points) || i < 0 {
		return nil
	}

	return &p.points[i]
}

// InsertAt inserts a Point at i along the path.
// Panics if index is out of range.
func (p *Path) InsertAt(index int, point *Point) *Path {
	if index > len(p.points) || index < 0 {
		panic(fmt.Sprintf("geo: insert index out of range, requested: %d, length: %d", index, len(p.points)))
	}

	if index == len(p.points) {
		p.points = append(p.points, *point)
		return p
	}

	p.points = append(p.points, Point{})
	copy(p.points[index+1:], p.points[index:])
	p.points[index] = *point

	return p
}

// RemoveAt removes a Point at i along the path.
// Panics if index is out of range.
func (p *Path) RemoveAt(index int) *Path {
	if index >= len(p.points) || index < 0 {
		panic(fmt.Sprintf("geo: remove index out of range, requested: %d, length: %d", index, len(p.points)))
	}

	p.points = append(p.points[:index], p.points[index+1:]...)
	return p
}

// Push appends a point to the end of the path.
func (p *Path) Push(point *Point) *Path {
	p.points = append(p.points, *point)
	return p
}

// Pop removes and returns the last point.
func (p *Path) Pop() *Point {
	if len(p.points) == 0 {
		return nil
	}

	x := p.points[len(p.points)-1]
	p.points = p.points[:len(p.points)-1]

	return &x
}

// Length returns the number of points in the path.
func (p *Path) Length() int {
	return len(p.points)
}

// Equals compares two paths. Returns true if lengths are the same
// and all points are Equal.
func (p *Path) Equals(path *Path) bool {
	if p.Length() != path.Length() {
		return false
	}

	for i, v := range p.points {
		if !v.Equals(&path.points[i]) {
			return false
		}
	}

	return true
}

// Clone returns a new copy of the path.
func (p *Path) Clone() *Path {
	points := make([]Point, len(p.points))
	copy(points, p.points)

	return &Path{
		points: points,
	}
}

// WriteOffFile writes an Object File Format representation of
// the points of the path to the writer provided. This is for viewing
// in MeshLab or something like that. You should close the
// writer yourself after this function returns.
// http://segeval.cs.princeton.edu/public/off_format.html
func (p *Path) WriteOffFile(w io.Writer, rgb ...[3]int) {
	r := 170
	g := 170
	b := 170

	if len(rgb) != 0 {
		r = rgb[0][0]
		g = rgb[0][1]
		b = rgb[0][2]
	}

	w.Write([]byte("OFF\n"))
	w.Write([]byte(fmt.Sprintf("%d %d 0\n", p.Length(), p.Length()-2)))

	for i := range p.points {
		w.Write([]byte(fmt.Sprintf("%f %f 0\n", p.points[i][0], p.points[i][1])))
	}

	for i := 0; i < len(p.points)-2; i++ {
		w.Write([]byte(fmt.Sprintf("3 %d %d %d %d %d %d\n", i, i+1, i+2, r, g, b)))
	}
}
