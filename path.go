package geo

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/paulmach/go.geojson"
)

// Path represents a set of points to be thought of as a polyline.
type Path []Point

// NewPath creates a new path.
func NewPath() Path {
	return NewPathPreallocate(0, 100)
}

// NewPathPreallocate creates a new path with points array of the given size.
func NewPathPreallocate(length, capacity int) Path {
	return Path(make([]Point, length, capacity))
}

// NewPathFromEncoding is the inverse of path.Encode. It takes a string encoding of a lat/lng path
// and returns the actual path it represents. Factor defaults to 1.0e5,
// the same used by Google for polyline encoding.
func NewPathFromEncoding(encoded string, factor ...int) Path {
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

			p = append(p, Point{float64(tempLatLng[1]) / f, float64(tempLatLng[0]) / f})
		}

		count++
	}

	return p
}

// NewPathFromXYData creates a path from a slice of [2]float64 values
// representing [horizontal, vertical] type data, for example lng/lat values from geojson.
func NewPathFromXYData(data [][2]float64) Path {
	p := NewPathPreallocate(0, len(data))
	for i := range data {
		p = append(p, Point{data[i][0], data[i][1]})
	}

	return p
}

// NewPathFromYXData creates a path from a slice of [2]float64 values
// representing [vertical, horizontal] type data, for example typical lat/lng data.
func NewPathFromYXData(data [][2]float64) Path {
	p := NewPathPreallocate(0, len(data))
	for i := range data {
		p = append(p, Point{data[i][1], data[i][0]})
	}

	return p
}

// NewPathFromXYSlice creates a path from a slice of []float64 values.
// The first two elements are taken to be horizontal and vertical components of each point respectively.
// The rest of the elements of the slice are ignored. Nil slices are skipped.
func NewPathFromXYSlice(data [][]float64) Path {
	p := NewPathPreallocate(0, len(data))
	for i := range data {
		if data[i] != nil && len(data[i]) >= 2 {
			p = append(p, Point{data[i][0], data[i][1]})
		}
	}

	return p
}

// NewPathFromYXSlice creates a path from a slice of []float64 values.
// The first two elements are taken to be vertical and horizontal components of each point respectively.
// The rest of the elements of the slice are ignored. Nil slices are skipped.
func NewPathFromYXSlice(data [][]float64) Path {
	p := NewPathPreallocate(0, len(data))
	for i := range data {
		if data[i] != nil && len(data[i]) >= 2 {
			p = append(p, Point{data[i][1], data[i][0]})
		}
	}

	return p
}

// Transform applies a given projection or inverse projection to all
// the points in the path.
func (p Path) Transform(projector Projector) Path {
	for i := range p {
		p[i] = projector(p[i])
	}

	return p
}

// Encode converts the path to a string using the Google Maps Polyline Encoding method.
// Factor defaults to 1.0e5, the same used by Google for polyline encoding.
func (p Path) Encode(factor ...int) string {
	f := 1.0e5
	if len(factor) != 0 {
		f = float64(factor[0])
	}

	var pLat int
	var pLng int

	var result bytes.Buffer
	scratch1 := make([]byte, 0, 50)
	scratch2 := make([]byte, 0, 50)

	for _, p := range p {
		lat5 := int(math.Floor(p.Lat()*f + 0.5))
		lng5 := int(math.Floor(p.Lng()*f + 0.5))

		deltaLat := lat5 - pLat
		deltaLng := lng5 - pLng

		pLat = lat5
		pLng = lng5

		result.Write(append(encodeSignedNumber(deltaLat, scratch1), encodeSignedNumber(deltaLng, scratch2)...))

		scratch1 = scratch1[:0]
		scratch2 = scratch2[:0]
	}

	return result.String()
}

func encodeSignedNumber(num int, result []byte) []byte {
	shiftedNum := num << 1

	if num < 0 {
		shiftedNum = ^shiftedNum
	}

	for shiftedNum >= 0x20 {
		result = append(result, byte(0x20|(shiftedNum&0x1f)+63))
		shiftedNum >>= 5
	}

	return append(result, byte(shiftedNum+63))
}

// Distance computes the total distance in the units of the points.
func (p Path) Distance() float64 {
	sum := 0.0

	loopTo := len(p) - 1
	for i := 0; i < loopTo; i++ {
		sum += p[i].DistanceFrom(p[i+1])
	}

	return sum
}

// GeoDistance computes the total distance using spherical geometry.
func (p Path) GeoDistance(haversine ...bool) float64 {
	yesgeo := yesHaversine(haversine)
	sum := 0.0

	loopTo := len(p) - 1
	for i := 0; i < loopTo; i++ {
		sum += p[i].GeoDistanceFrom(p[i+1], yesgeo)
	}

	return sum
}

// DistanceFrom computes an O(n) distance from the path. Loops over every
// subline to find the minimum distance.
func (p Path) DistanceFrom(point Point) float64 {
	return math.Sqrt(p.SquaredDistanceFrom(point))
}

// SquaredDistanceFrom computes an O(n) minimum squared distance from the path.
// Loops over every subline to find the minimum distance.
func (p Path) SquaredDistanceFrom(point Point) float64 {
	dist := math.Inf(1)

	l := Line{}
	loopTo := len(p) - 1
	for i := 0; i < loopTo; i++ {
		l.a = p[i]
		l.b = p[i+1]
		dist = math.Min(l.SquaredDistanceFrom(point), dist)
	}

	return dist
}

// DirectionAt computes the direction of the path at the given index.
// Uses the line between the two surrounding points to get the direction,
// or just the first two, or last two if at the start or end, respectively.
// Assumes the path is in a conformal projection.
// The units are radians from the positive x-axis. Range same as math.Atan2, [-Pi, Pi]
// Returns INF for single point paths.
func (p Path) DirectionAt(index int) float64 {
	if index >= len(p) || index < 0 {
		panic("TODO")
		// panic(fmt.Sprintf("geo: direction at index out of range, requested: %d, length: %d", index, len(p.PointSet)))
	}

	if len(p) == 1 {
		return math.Inf(1)
	}

	var diff Vector
	if index == 0 {
		diff = p[1].Subtract(p[0])
	} else if index >= len(p)-1 {
		length := len(p)
		diff = p[length-1].Subtract(p[length-2])
	} else {
		diff = p[index+1].Subtract(p[index-1])
	}

	return math.Atan2(diff.Y(), diff.X())
}

// Measure computes the distance along this path to the point nearest the given point.
func (p Path) Measure(point Point) float64 {
	minDistance := math.Inf(1)
	measure := math.Inf(-1)
	sum := 0.0

	seg := Line{}
	for i := 0; i < len(p)-1; i++ {
		seg.a = p[i]
		seg.b = p[i+1]
		distanceToLine := seg.SquaredDistanceFrom(point)
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
func (p Path) Project(point Point) float64 {
	return p.Measure(point) / p.Distance()
}

// Intersection calls IntersectionPath or IntersectionLine depending on the
// type of the provided geometry.
// TODO: have this receive an Intersectable interface.
func (p Path) Intersection(geometry interface{}) ([]Point, [][2]int) {
	var points []Point
	var segments [][2]int

	switch g := geometry.(type) {
	case Line:
		points, segments = p.IntersectionLine(g)
	case *Line:
		points, segments = p.IntersectionLine(*g)
	case Path:
		points, segments = p.IntersectionPath(g)
	case *Path:
		points, segments = p.IntersectionPath(*g)
	default:
		panic("can only determine intersection with lines and paths")
	}

	return points, segments
}

// IntersectionPath returns a slice of points and a slice of tuples [i, j] where i is the segment
// in the parent path and j is the segment in the given path that intersect to form the given point.
// Slices will be empty if there is no intersection.
func (p Path) IntersectionPath(path Path) ([]Point, [][2]int) {
	// TODO: done some sort of line sweep here if p.Length() is big enough
	var points []Point
	var indexes [][2]int

	for i := 0; i < len(p)-1; i++ {
		pLine := NewLine(p[i], p[i+1])

		for j := 0; j < len(path)-1; j++ {
			pathLine := NewLine(path[j], path[j+1])

			if point, err := pLine.Intersection(pathLine); err == nil {
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
func (p Path) IntersectionLine(line Line) ([]Point, [][2]int) {
	var points []Point
	var indexes [][2]int

	for i := 0; i < len(p)-1; i++ {
		pTest := NewLine(p[i], p[i+1])
		if point, err := pTest.Intersection(line); err == nil {
			points = append(points, point)
			indexes = append(indexes, [2]int{i, 0})
		}
	}

	return points, indexes
}

// Intersects can take a line or a path to determine if there is an intersection.
// TODO: I would love this to accept an intersecter interface.
func (p Path) Intersects(geometry interface{}) bool {
	var result bool

	switch g := geometry.(type) {
	case Line:
		result = p.IntersectsLine(g)
	case *Line:
		result = p.IntersectsLine(*g)
	case Path:
		result = p.IntersectsPath(g)
	case *Path:
		result = p.IntersectsPath(*g)
	default:
		panic("can only determine intersection with lines and paths")
	}

	return result
}

// IntersectsPath takes a Path and checks if it intersects with the path.
func (p Path) IntersectsPath(path Path) bool {
	// TODO: done some sort of line sweep here if p.Length() is big enough
	for i := 0; i < len(p)-1; i++ {
		pLine := NewLine(p[i], p[i+1])

		for j := 0; j < len(path)-1; j++ {
			pathLine := NewLine(path[j], path[j+1])

			if pLine.Intersects(pathLine) {
				return true
			}
		}
	}

	return false
}

// IntersectsLine takes a Line and checks if it intersects with the path.
func (p Path) IntersectsLine(line Line) bool {
	for i := 0; i < len(p)-1; i++ {
		pTest := NewLine(p[i], p[i+1])
		if pTest.Intersects(line) {
			return true
		}
	}

	return false
}

// Bound returns a bound around the path. Uses rectangular coordinates.
func (p Path) Bound() Bound {
	if len(p) == 0 {
		return NewBound(0, 0, 0, 0)
	}

	minX := math.Inf(1)
	minY := math.Inf(1)

	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	for _, v := range p {
		minX = math.Min(minX, v.X())
		minY = math.Min(minY, v.Y())

		maxX = math.Max(maxX, v.X())
		maxY = math.Max(maxY, v.Y())
	}

	return NewBound(maxX, minX, maxY, minY)
}

// Equal compares two paths. Returns true if lengths are the same
// and all points are Equal.
func (p Path) Equal(path Path) bool {
	if len(p) != len(path) {
		return false
	}

	for i := range p {
		if !p[i].Equal(path[i]) {
			return false
		}
	}

	return true
}

// Clone returns a new copy of the path.
func (p Path) Clone() Path {
	ps := PointSet(p)
	return Path(ps.Clone())
}

// ToGeoJSON creates a new geojson feature with a linestring geometry
// containing all the points.
func (p Path) ToGeoJSON() *geojson.Geometry {
	coords := make([][]float64, 0, len(p))

	for _, point := range p {
		coords = append(coords, []float64{point[0], point[1]})
	}

	return geojson.NewLineStringGeometry(coords)
}

// ToWKT returns the path in WKT format, eg. LINESTRING(30 10,10 30,40 40)
// For empty paths the result will be 'EMPTY'.
func (p Path) ToWKT() string {
	return p.String()
}

// String returns a string representation of the path.
// The format is WKT, e.g. LINESTRING(30 10,10 30,40 40)
// For empty paths the result will be 'EMPTY'.
func (p Path) String() string {
	if len(p) == 0 {
		return "EMPTY"
	}

	buff := bytes.NewBuffer(nil)
	fmt.Fprintf(buff, "LINESTRING(%g %g", p[0][0], p[0][1])

	for i := 1; i < len(p); i++ {
		fmt.Fprintf(buff, ",%g %g", p[i][0], p[i][1])
	}

	buff.Write([]byte(")"))
	return buff.String()
}

// WriteOffFile writes an Object File Format representation of
// the points of the path to the writer provided. This is for viewing
// in MeshLab or something like that. You should close the
// writer yourself after this function returns.
// http://segeval.cs.princeton.edu/public/off_format.html
func (p Path) WriteOffFile(w io.Writer, rgb ...[3]int) {
	r := 170
	g := 170
	b := 170

	if len(rgb) != 0 {
		r = rgb[0][0]
		g = rgb[0][1]
		b = rgb[0][2]
	}

	w.Write([]byte("OFF\n"))
	w.Write([]byte(fmt.Sprintf("%d %d 0\n", len(p), len(p)-2)))

	for i := range p {
		w.Write([]byte(fmt.Sprintf("%f %f 0\n", p[i][0], p[i][1])))
	}

	for i := 0; i < len(p)-2; i++ {
		w.Write([]byte(fmt.Sprintf("3 %d %d %d %d %d %d\n", i, i+1, i+2, r, g, b)))
	}
}
