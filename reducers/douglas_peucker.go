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
func NewDouglasPeucker(threshold float64) DouglasPeuckerReducer {
	return DouglasPeuckerReducer{
		Threshold: threshold,
	}
}

// Reduce runs the DouglasPeucker using the threshold of the DouglasPeuckerReducer.
func (r DouglasPeuckerReducer) Reduce(path geo.Path) geo.Path {
	return DouglasPeucker(path, r.Threshold)
}

// GeoReduce runs the DouglasPeucker on a lng/lat path.
// The threshold is expected to be in meters.
func (r DouglasPeuckerReducer) GeoReduce(path geo.Path) geo.Path {
	factor := geo.MercatorScaleFactor(path.Bound().Center().Lat())
	path.Transform(geo.Mercator.Project)
	reduced := DouglasPeucker(path, r.Threshold*factor)

	return reduced.Transform(geo.Mercator.Inverse)
}

// DouglasPeucker simplifies the path using the Douglas Peucker method.
// Returns a new path and DOES NOT modify the original.
func DouglasPeucker(path geo.Path, threshold float64) geo.Path {
	if len(path) <= 2 {
		return path.Clone()
	}

	mask := make([]byte, len(path))
	mask[0] = 1
	mask[len(path)-1] = 1

	found := dpWorker(path, threshold, mask)
	newPath := geo.NewPathPreallocate(0, found)

	for i, v := range mask {
		if v == 1 {
			newPath = append(newPath, path[i])
		}
	}

	return newPath
}

// DouglasPeuckerIndexMap is similar to DouglasPeucker but returns an array that maps
// each new path index to its original path index.
// Returns a new path and DOES NOT modify the original.
func DouglasPeuckerIndexMap(path geo.Path, threshold float64) (reduced geo.Path, indexMap []int) {
	if len(path) == 0 {
		return path.Clone(), []int{}
	}

	if len(path) == 1 {
		return path.Clone(), []int{0}
	}

	if len(path) == 2 {
		return path.Clone(), []int{0, 1}
	}

	mask := make([]byte, len(path))
	mask[0] = 1
	mask[len(path)-1] = 1

	found := dpWorker(path, threshold, mask)

	newPath := geo.NewPathPreallocate(0, found)
	for i, v := range mask {
		if v == 1 {
			newPath = append(newPath, path[i])
			indexMap = append(indexMap, i)
		}
	}

	return newPath, indexMap
}

// dpWorker does the recursive threshold checks.
// Using a stack array with a stackLength variable resulted in 4x speed improvement
// over calling the function recursively.
func dpWorker(path geo.Path, threshold float64, mask []byte) int {

	found := 0

	var stack []int
	stack = append(stack, 0, len(path)-1)

	for len(stack) > 0 {
		start := stack[len(stack)-2]
		end := stack[len(stack)-1]

		// modify the line in place
		l := geo.NewLine(
			geo.NewPoint(path[start][0], path[start][1]),
			geo.NewPoint(path[end][0], path[end][1]),
		)

		maxDist := 0.0
		maxIndex := 0
		for i := start + 1; i < end; i++ {
			dist := l.SquaredDistanceFrom(path[i])

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
