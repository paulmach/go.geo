package reducers

import (
	"github.com/paulmach/go.geo"
)

// A DouglasPeuckerReducer wraps the DouglasPeucker function
// to fulfil the geo.Reducer interface.
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

// DouglasPeucker simplifies the path using the Douglas Peucker method.
// Returns a new path and DOES NOT modify the original.
func DouglasPeucker(path *geo.Path, threshold float64) *geo.Path {
	if path.Length() <= 2 {
		return path.Clone()
	}

	mask := make([]byte, path.Length())
	mask[0] = 1
	mask[path.Length()-1] = 1

	dpWorker(path, 0, path.Length()-1, threshold, mask)

	count := 0
	points := path.Points()

	for i, v := range mask {
		if v == 1 {
			points[count] = points[i]
			count++
		}
	}

	points = points[:count]
	return path.SetPoints(points)
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

	dpWorker(path, 0, path.Length()-1, threshold, mask)

	originalPoints := path.Points()
	var points []geo.Point

	for i, v := range mask {
		if v == 1 {
			points = append(points, originalPoints[i])
			indexMap = append(indexMap, i)
		}
	}

	reduced = &geo.Path{}
	return reduced.SetPoints(points), indexMap
}

func dpWorker(path *geo.Path, start, end int, threshold float64, mask []byte) {
	if end-start <= 1 {
		return
	}

	l := geo.NewLine(path.GetAt(start), path.GetAt(end))

	maxDist := 0.0
	maxIndex := start + 1
	for i := start + 1; i < end; i++ {
		dist := l.DistanceFrom(path.GetAt(i))

		if dist >= maxDist {
			maxDist = dist
			maxIndex = i
		}
	}

	if maxDist > threshold {
		mask[maxIndex] = 1

		dpWorker(path, start, maxIndex, threshold, mask)
		dpWorker(path, maxIndex, end, threshold, mask)
	}
}
