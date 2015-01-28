package geo

import (
	"fmt"
	"math"
)

// A PointSet represents a set of points in the 2D Eucledian or Cartesian plane.
type PointSet []Point

// NewPointSetPreallocate simply creates a new point set with points array of the given size.
func NewPointSet() *PointSet {
	return &PointSet{}
}

// NewPointSetPreallocate simply creates a new point set with points array of the given size.
func NewPointSetPreallocate(length, capacity int) *PointSet {
	if length > capacity {
		capacity = length
	}

	ps := make([]Point, length, capacity)
	p := PointSet(ps)
	return &p
}

// Clone returns a new copy of the point set.
func (ps PointSet) Clone() PointSet {
	points := make([]Point, len(ps))
	copy(points, ps)

	return points
}

// Centroid returns the average latitude and longitude coordinate of the point set
func (ps PointSet) GeoCentroid() *Point {
	averageLat := 0.0
	averageLng := 0.0
	numPoints := float64(len(ps))
	for _, point := range ps {
		averageLat += point.Lat()
		averageLng += point.Lng()
	}
	return &Point{averageLng / numPoints, averageLat / numPoints}
}

// GeoDistanceFrom returns the minimum geo distance from the point set
func (ps PointSet) GeoDistanceFrom(point *Point) float64 {
	dist := math.Inf(1)

	loopTo := len(ps) - 1
	for i := 0; i < loopTo; i++ {
		dist = math.Min(ps[i].GeoDistanceFrom(point), dist)
	}

	return dist
}

// SetAt updates a position at i in the point set
func (ps *PointSet) SetAt(index int, point *Point) *PointSet {
	deref := *ps
	if index >= len(deref) || index < 0 {
		panic(fmt.Sprintf("geo: set index out of range, requested: %d, length: %d", index, len(deref)))
	}
	deref[index] = *point
	*ps = deref
	return ps
}

// GetAt returns the pointer to the Point in the page.
// This function is good for modifying values in place.
// Returns nil if index is out of range.
func (ps *PointSet) GetAt(i int) *Point {
	deref := *ps
	if i >= len(deref) || i < 0 {
		return nil
	}

	return &deref[i]
}

// InsertAt inserts a Point at i in the point set.
// Panics if index is out of range.
func (ps *PointSet) InsertAt(index int, point *Point) *PointSet {
	deref := *ps
	if index > len(deref) || index < 0 {
		panic(fmt.Sprintf("geo: insert index out of range, requested: %d, length: %d", index, len(deref)))
	}

	if index == len(deref) {
		deref = append(deref, *point)
		*ps = deref
		return ps
	}

	deref = append(deref, Point{})
	copy(deref[index+1:], deref[index:])
	deref[index] = *point
	*ps = deref
	return ps
}

// RemoveAt removes a Point at i in the point set.
// Panics if index is out of range.
func (ps *PointSet) RemoveAt(index int) *PointSet {
	deref := *ps
	if index >= len(deref) || index < 0 {
		panic(fmt.Sprintf("geo: remove index out of range, requested: %d, length: %d", index, len(deref)))
	}

	deref = append(deref[:index], deref[index+1:]...)
	*ps = deref
	return ps
}

// Push appends a point to the end of the point set.
func (ps *PointSet) Push(point *Point) *PointSet {
	*ps = append(*ps, *point)
	return ps
}

// Pop removes and returns the last point in the point set
func (ps *PointSet) Pop() *Point {
	deref := *ps
	if len(deref) == 0 {
		return nil
	}

	x := deref[len(deref)-1]
	*ps = deref[:len(deref)-1]

	return &x
}

//SetPoints sets the points in the point set
func (ps *PointSet) SetPoints(points []Point) *PointSet {
	*ps = points
	return ps
}

// Length returns the number of points in the point set.
func (ps PointSet) Length() int {
	return len(ps)
}

// Equals compares two point sets. Returns true if lengths are the same
// and all points are Equal
func (ps *PointSet) Equals(pointSet *PointSet) bool {
	if (*ps).Length() != (*pointSet).Length() {
		return false
	}

	for i, v := range *ps {
		if !v.Equals((*pointSet).GetAt(i)) {
			return false
		}
	}

	return true
}
