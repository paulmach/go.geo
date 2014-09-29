package reducers

import (
	"github.com/paulmach/go.geo"
)

// A DouglasPeuckerReducer wraps the DouglasPeucker function
// to fulfill the geo.Reducer and geo.GeoReducer interfaces.
type DouglasPeuckerReducer struct {
	Threshold float64
}

// NewDouglasPeucker creates a new DouglasPeuckerReducer.
func NewDouglasPeucker(threshold float64) *DouglasPeuckerReducer {
	return &DouglasPeuckerReducer{
		Threshold: threshold,
	}
}

// Reduce runs the DouglasPeucker using the threshold of the DouglasPeuckerReducer.
func (r DouglasPeuckerReducer) Reduce(path *geo.Path) *geo.Path {
	return DouglasPeucker(path, r.Threshold)
}

// GeoReduce runs the DouglasPeucker on a lng/lat path.
// The threshold is expected to be in meters.
func (r DouglasPeuckerReducer) GeoReduce(path *geo.Path) *geo.Path {
	factor := geo.MercatorScaleFactor(path.Bound().Center().Lat())
	path.Transform(geo.Mercator.Project)
	reduced := DouglasPeucker(path, r.Threshold*factor)

	return reduced.Transform(geo.Mercator.Inverse)
}

// DouglasPeucker simplifies the path using the Douglas Peucker method.
// Returns a new path and DOES NOT modify the original.
func DouglasPeucker(path *geo.Path, threshold float64) *geo.Path {
	if path.Length() <= 2 {
		return path.Clone()
	}

	mask := make([]byte, path.Length())
	mask[0] = 1
	mask[path.Length()-1] = 1

	points := path.Points()

	found := dpWorker(points, threshold, mask)
	newPoints := make([]geo.Point, 0, found)

	for i, v := range mask {
		if v == 1 {
			newPoints = append(newPoints, points[i])
		}
	}

	return (&geo.Path{}).SetPoints(newPoints)
}

// DouglasPeuckerIndexMap is similar to DouglasPeucker but returns an array that maps
// each new path index to its original path index.
// Returns a new path and DOES NOT modify the original.
func DouglasPeuckerIndexMap(path *geo.Path, threshold float64) (reduced *geo.Path, indexMap []int) {
	if path.Length() == 0 {
		return path.Clone(), []int{}
	}

	if path.Length() == 1 {
		return path.Clone(), []int{0}
	}

	if path.Length() == 2 {
		return path.Clone(), []int{0, 1}
	}

	mask := make([]byte, path.Length())
	mask[0] = 1
	mask[path.Length()-1] = 1

	originalPoints := path.Points()

	found := dpWorker(originalPoints, threshold, mask)

	points := make([]geo.Point, 0, found)
	for i, v := range mask {
		if v == 1 {
			points = append(points, originalPoints[i])
			indexMap = append(indexMap, i)
		}
	}

	reduced = &geo.Path{}
	return reduced.SetPoints(points), indexMap
}

// dpWorker does the recursive threshold checks.
// Using a stack array with a stackLength variable resulted in 4x speed improvement
// over calling the function recursively.
func dpWorker(points []geo.Point, threshold float64, mask []byte) int {

	found := 0

	var stack []int
	stack = append(stack, 0, len(points)-1)

	l := &geo.Line{}
	for len(stack) > 0 {
		start := stack[len(stack)-2]
		end := stack[len(stack)-1]

		// modify the line in place
		a := l.A()
		a[0], a[1] = points[start][0], points[start][1]

		b := l.B()
		b[0], b[1] = points[end][0], points[end][1]

		maxDist := 0.0
		maxIndex := 0
		for i := start + 1; i < end; i++ {
			dist := l.SquaredDistanceFrom(&points[i])

			if dist > maxDist {
				maxDist = dist
				maxIndex = i
			}
		}

		if maxDist > threshold*threshold {
			found++
			mask[maxIndex] = 1

			stack[len(stack)-1] = maxIndex
			stack = append(stack, maxIndex, end)
		} else {
			stack = stack[:len(stack)-2]
		}
	}

	return found
}
