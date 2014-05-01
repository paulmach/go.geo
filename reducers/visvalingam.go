package reducers

import (
	"container/heap"
	"container/list"
	"math"

	"github.com/paulmach/go.geo"
)

// VisvalingamThreshold does the Visvalingam-Whyatt algorithm removing
// triangles whose area is below the threshold. This function is here to simplify the interface.
// Returns a new path and DOES NOT modify the original.
func VisvalingamThreshold(path *geo.Path, threshold float64) *geo.Path {
	return Visvalingam(path, threshold, 0)
}

// VisvalingamKeep does the Visvalingam-Whyatt algorithm removing
// triangles of minimum area until we're down to toKeep number of points.
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

	if path.Length() <= 1 {
		return path.Clone()
	}

	// since our triangleSquareDistanceNormalArea function
	// doesn't do the extra 1/4 sqrt we do the opposite to the threshold
	threshold = 16 * threshold * threshold

	numPoints := path.Length()
	points := path.Points()

	removed := 0

	// build the initial linked list and priority queue
	pointList := list.New()
	queue := &visQueue{}
	heap.Init(queue)

	pointList.PushBack(&visItem{
		point: &points[0],
		loc:   0,
		area:  -1,
		index: -1,
	})

	d1 := sqDistance(&points[0], &points[1])
	for i := 1; i < numPoints-1; i++ {
		d2 := sqDistance(&points[i], &points[i+1])
		d3 := sqDistance(&points[i-1], &points[i+1])

		item := &visItem{
			point: &points[i],
			area:  triangleSquareDistanceNormalArea(d1, d2, d3),
			index: i - 1,
			loc:   i,
		}

		// add the item to the end of the linked list
		// also push that item on the heap/priority queue
		heap.Push(queue, pointList.PushBack(item))

		d1 = d2
	}

	pointList.PushBack(&visItem{
		point: &points[numPoints-1],
		area:  -2,
		index: -2,
		loc:   numPoints - 1,
	})

	// run through the reduction process
	for queue.Len() > 0 {
		element := heap.Pop(queue).(*list.Element)
		item := element.Value.(*visItem)

		if item.area > threshold || numPoints-removed <= minPointsToKeep {
			break
		}

		next := element.Next()
		nextnext := next.Next()
		prev := element.Prev()
		prevprev := prev.Prev()

		// remove current element from list
		pointList.Remove(element)
		removed++

		// figure out the new area of the previous element
		if prevprev != nil {
			area := trianglePointNormalArea(
				prevprev.Value.(*visItem).point,
				prev.Value.(*visItem).point,
				next.Value.(*visItem).point)

			area = math.Max(area, item.area)
			queue.update(prev, area)
		}

		if nextnext != nil {
			area := trianglePointNormalArea(
				prev.Value.(*visItem).point,
				next.Value.(*visItem).point,
				nextnext.Value.(*visItem).point)

			area = math.Max(area, item.area)
			queue.update(next, area)
		}
	}

	newPoints := make([]geo.Point, 0, queue.Len()+2)
	for e := pointList.Front(); e != nil; e = e.Next() {
		newPoints = append(newPoints, *(e.Value.(*visItem).point))
	}

	reduced := &geo.Path{}
	return reduced.SetPoints(newPoints)
}

// Stuff to create the priority queue and the related linked list.
//
// The priority queue is made up of elements from a linked of the points.
// This allows you to access the previous and next point, and remove points simply.

type visItem struct {
	point *geo.Point
	area  float64
	loc   int
	index int
}

type visQueue []*list.Element

func (q visQueue) Len() int { return len(q) }

func (q visQueue) Less(i, j int) bool {
	return q[i].Value.(*visItem).area < q[j].Value.(*visItem).area
}

func (q visQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].Value.(*visItem).index = i
	q[j].Value.(*visItem).index = j
}

func (q *visQueue) Push(x interface{}) {
	item := x.(*list.Element)
	item.Value.(*visItem).index = q.Len()
	*q = append(*q, item)
}

func (q *visQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.Value.(*visItem).index = -1
	*q = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (q *visQueue) update(element *list.Element, area float64) {
	heap.Remove(q, element.Value.(*visItem).index)
	element.Value.(*visItem).area = area
	heap.Push(q, element)
}

func sqDistance(a, b *geo.Point) float64 {
	return (a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1])
}

func trianglePointNormalArea(a, b, c *geo.Point) float64 {
	return triangleSquareDistanceNormalArea(sqDistance(a, b), sqDistance(b, c), sqDistance(a, c))
}

// triangleSquareDistanceNormalArea returns the triangle area  =  (4 * area)^2
func triangleSquareDistanceNormalArea(a, b, c float64) float64 {
	return 2*(a*b+b*c+c*a) - a*a - b*b - c*c
}
