package reducers

import (
	"math"

	"github.com/paulmach/go.geo"
)

// A VisvalingamReducer wraps the Visvalingam function
// to fulfill the geo.Reducer and geo.GeoReducer interfaces.
type VisvalingamReducer struct {
	Threshold float64
	ToKeep    int
}

// NewVisvalingamReducer creates a new VisvalingamReducer.
func NewVisvalingamReducer(threshold float64, minPointsToKeep int) *VisvalingamReducer {
	return &VisvalingamReducer{
		Threshold: threshold,
		ToKeep:    minPointsToKeep,
	}
}

// Reduce runs the Visvalingam reduction using the values of the Visvalingam.
func (r VisvalingamReducer) Reduce(path *geo.Path) *geo.Path {
	return Visvalingam(path, r.Threshold, r.ToKeep)
}

// GeoReduce runs the Visvalingam reduction on a lng/lat path.
// The threshold is expected to be in meters squared.
func (r VisvalingamReducer) GeoReduce(path *geo.Path) *geo.Path {
	factor := geo.MercatorScaleFactor(path.Bound().Center().Lat())
	path.Transform(geo.Mercator.Project)

	reduced := Visvalingam(path, r.Threshold*factor*factor, r.ToKeep)
	return reduced.Transform(geo.Mercator.Inverse)
}

// VisvalingamThreshold runs the Visvalingam-Whyatt algorithm removing
// triangles whose area is below the threshold. This function is here to simplify the interface.
// Returns a new path and DOES NOT modify the original.
func VisvalingamThreshold(path *geo.Path, threshold float64) *geo.Path {
	return Visvalingam(path, threshold, 0)
}

// VisvalingamKeep runs the Visvalingam-Whyatt algorithm removing
// triangles of minimum area until we're down to `toKeep` number of points.
// Returns a new path and DOES NOT modify the original.
func VisvalingamKeep(path *geo.Path, toKeep int) *geo.Path {
	return Visvalingam(path, math.MaxFloat64, toKeep)
}

// Visvalingam computes the Visvalingam-Whyatt on the polyline.
// Returns a new path and DOES NOT modify the original.
//
// Threshold is the max triangle area to keep around, ie. remove all triangles below this threshold.
// minPointsToKeep is the minimum number of points in the line.
//
// To just use the threshold, set minPointsToKeep to zero
// To just use minPointsToKeep, set the threshold to something big like math.MaxFloat64
//
// http://bost.ocks.org/mike/simplify/
func Visvalingam(path *geo.Path, threshold float64, minPointsToKeep int) *geo.Path {
	if threshold < 0 {
		panic("threshold must be >= 0")
	}

	if path.Length() <= minPointsToKeep {
		return path.Clone()
	}

	if path.Length() <= 2 {
		return path.Clone()
	}

	// edge cases checked, get on with it
	threshold *= 2 // triangle area is doubled to save the multiply :)
	removed := 0

	points := path.Points()
	numPoints := len(points)

	// build the initial minheap linked list.
	heap := minHeap(make([]*visItem, 0, numPoints))

	linkedListStart := &visItem{
		area:       math.Inf(1),
		pointIndex: 0,
	}
	heap.Push(linkedListStart)

	// internal path items
	previous := linkedListStart
	for i := 1; i < numPoints-1; i++ {
		current := &visItem{
			area:       doubleTriangleArea(&points[i-1], &points[i], &points[i+1]),
			pointIndex: i,
			previous:   previous,
		}

		heap.Push(current)
		previous.next = current
		previous = current
	}

	// final item
	endItem := &visItem{
		area:       math.Inf(1),
		pointIndex: numPoints - 1,
		previous:   previous,
	}
	previous.next = endItem
	heap.Push(endItem)

	// run through the reduction process
	for len(heap) > 0 {
		current := heap.Pop()

		if current.area > threshold || numPoints-removed <= minPointsToKeep {
			break
		}

		next := current.next
		previous := current.previous

		// remove current element from linked list
		previous.next = current.next
		next.previous = current.previous
		removed++

		// figure out the new areas
		if previous.previous != nil {
			area := doubleTriangleArea(
				&points[previous.previous.pointIndex],
				&points[previous.pointIndex],
				&points[next.pointIndex],
			)

			area = math.Max(area, current.area)
			heap.Update(previous, area)
		}

		if next.next != nil {
			area := doubleTriangleArea(
				&points[previous.pointIndex],
				&points[next.pointIndex],
				&points[next.next.pointIndex],
			)

			area = math.Max(area, current.area)
			heap.Update(next, area)
		}
	}

	item := linkedListStart
	newPoints := make([]geo.Point, 0, len(heap)+2)

	for item != nil {
		newPoints = append(newPoints, points[item.pointIndex])
		item = item.next
	}

	reduced := &geo.Path{}
	return reduced.SetPoints(newPoints)
}

// Stuff to create the priority queue, or min heap.
// Rewriting it here, vs using the std lib, resulted in a 10x performance bump!
type minHeap []*visItem

type visItem struct {
	area       float64 // triangle area
	pointIndex int     // index of point in original path

	// to keep a virtual linked list to help rebuild the triangle areas as we remove points.
	next     *visItem
	previous *visItem

	index int // interal index in heap, for removal and update
}

func (h *minHeap) Push(item *visItem) {
	item.index = len(*h)
	*h = append(*h, item)
	h.up(item.index)
}

func (h *minHeap) Pop() *visItem {
	removed := (*h)[0]
	lastItem := (*h)[len(*h)-1]
	(*h) = (*h)[:len(*h)-1]

	if len(*h) > 0 {
		lastItem.index = 0
		(*h)[0] = lastItem
		h.down(0)
	}

	return removed
}

func (h minHeap) Update(item *visItem, area float64) {
	if item.area > area {
		// area got smaller
		item.area = area
		h.up(item.index)
	} else {
		// area got larger
		item.area = area
		h.down(item.index)
	}
}

func (h *minHeap) Remove(item *visItem) {
	i := item.index

	lastItem := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]

	if i != len(*h) {
		lastItem.index = i
		(*h)[i] = lastItem

		if lastItem.area < item.area {
			h.up(i)
		} else {
			h.down(i)
		}
	}
}

func (h minHeap) up(i int) {
	object := h[i]
	for i > 0 {
		up := ((i + 1) >> 1) - 1
		parent := h[up]

		if parent.area <= object.area {
			// parent is smaller so we're done fixing up the heap.
			break
		}

		// swap nodes
		parent.index = i
		h[i] = parent

		object.index = up
		h[up] = object

		i = up
	}
}

func (h minHeap) down(i int) {
	object := h[i]
	for {
		right := (i + 1) << 1
		left := right - 1

		down := i
		child := h[down]

		// swap with smallest child
		if left < len(h) && h[left].area < child.area {
			down = left
			child = h[down]
		}

		if right < len(h) && h[right].area < child.area {
			down = right
			child = h[down]
		}

		// non smaller, so quit
		if down == i {
			break
		}

		// swap the nodes
		child.index = i
		h[child.index] = child

		object.index = down
		h[down] = object

		i = down
	}
}

func doubleTriangleArea(a, b, c *geo.Point) float64 {
	area := (b[0]-a[0])*(c[1]-a[1]) - (b[1]-a[1])*(c[0]-a[0])
	return math.Abs(area)
}
