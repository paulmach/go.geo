// Package quadtree implements a quadtree using rectangular partitions.
// Each point exists in a unique node; if multiple points are in the same position,
// some points may be stored on internal nodes rather than leaf nodes.
// This implementation is based heavily off of the d3 implementation:
// https://github.com/mbostock/d3/wiki/Quadtree-Geom
package quadtree

import (
	"errors"
	"math"

	"github.com/paulmach/go.geo"
)

var (
	// ErrPointOutsideOfBounds is returned when trying to add a point
	// to a quad tree and the point is outside the bounds used to create the tree.
	ErrPointOutsideOfBounds = errors.New("quadtree: point outside of bounds")
)

// Quadtree implements a two-dimensional recursive spatial subdivision
// of geo.Pointers. This implementation uses rectangular partitions.
type Quadtree struct {
	// Threshold indicates the limit of how deep the quadtree can go.
	// Points closer than this will essentially be put in the same leaf node and stop recursion.
	// The correct value depends on the use case. The default is computed
	// off the bounds to keep the tree at most 12 levels deep. So points that
	// are closer than 1/4096 * max(bound.width, bound.height) will essentially be
	// stored in the same leaf node. For optimal tree performance you want this to happen
	// sometimes but not very often.
	Threshold float64

	bound     *geo.Bound
	root      *node
	freeNodes []node
	freeIndex int
}

// node represents a node of the quad tree. Each node stores a Value
// and has links to its 4 children
type node struct {
	children [4]*node
	internal bool
	pointer  geo.Pointer
}

// New creates a new quadtree for the given bound. Added points
// must be within this bound.
func New(bound *geo.Bound) *Quadtree {
	return &Quadtree{
		Threshold: math.Max(bound.Width(), bound.Height()) / float64(1<<12),
		bound:     bound,
	}
}

// NewFromPointSet creates a quadtree from a pointset.
// Copies the points into the quad tree. Modifying the points later
// will invalidate the quad tree and lead to unexpected result.
func NewFromPointSet(set *geo.PointSet) *Quadtree {
	q := New(set.Bound())
	q.freeNodes = make([]node, set.Length(), set.Length())

	ps := []geo.Point(*set)
	for i := range ps {
		q.Insert(&ps[i])
	}

	return q
}

// NewFromPointers creates a quadtree from a set of pointers.
func NewFromPointers(points []geo.Pointer) *Quadtree {
	if len(points) == 0 {
		// This is kind of meaningless but is what will happen
		// if using an empty pointset above.
		return New(geo.NewBound(0, 0, 0, 0))
	}

	b := geo.NewBoundFromPoints(points[0].Point(), points[0].Point())
	for _, p := range points {
		b.Extend(p.Point())
	}

	q := New(b)
	q.freeNodes = make([]node, len(points), len(points))

	for _, p := range points {
		q.Insert(p)
	}

	return q
}

// Bound returns the bounds used for the quad tree.
func (q *Quadtree) Bound() *geo.Bound {
	return q.bound
}

// Insert puts an object into the quad tree, must be within the quadtree bounds.
// If the pointer returns nil, the point will be ignored.
// This function is not thread-safe, ie. multiple goroutines cannot insert into
// a single quadtree.
func (q *Quadtree) Insert(p geo.Pointer) error {
	if p == nil {
		return nil
	}

	point := p.Point()
	if point == nil {
		return nil
	}

	if !q.bound.Contains(point) {
		return ErrPointOutsideOfBounds
	}

	if q.root == nil {
		q.root = &node{}
	}

	q.insert(q.root, p,
		q.bound.Left(), q.bound.Right(),
		q.bound.Bottom(), q.bound.Top(),
	)

	return nil
}

// nextNode returns the next node from a preallocated list.
// This resulted in about 15% improvement in quadtree creation.
func (q *Quadtree) nextNode() *node {
	if l := len(q.freeNodes); q.freeIndex >= l {
		// Exponentially decrease the preallocation size.
		// On a handful of tests, number of nodes was about 1.5 times pointers.
		l /= 2

		// min size of the preallocation. I think this could be bigger as it's
		// not that much memory overhead. Optimizing this more would need
		// to be use case specific.
		if l < 25 {
			l = 25
		}

		q.freeNodes = make([]node, l, l)
		q.freeIndex = 0
	}

	n := &q.freeNodes[q.freeIndex]
	q.freeIndex++
	return n
}

func (q *Quadtree) insert(n *node, p geo.Pointer, left, right, bottom, top float64) {
	point := p.Point()
	if n.internal {
		i := 0

		// figure which child of this internal node the point is in.
		if cy := (bottom + top) / 2.0; point.Y() <= cy {
			top = cy
			i = 2
		} else {
			bottom = cy
		}

		if cx := (left + right) / 2.0; point.X() >= cx {
			left = cx
			i++
		} else {
			right = cx
		}

		if n.children[i] == nil {
			// child not yet created so automatically add the pointer to it and return.
			n.children[i] = q.nextNode()
			n.children[i].pointer = p
			return
		}

		// proceed down to the child to see if it's a leaf yet and we can add the pointer there.
		q.insert(n.children[i], p, left, right, bottom, top)
		return
	}

	if n.pointer == nil {
		// leaf without a pointer. I believe this only happens for the first pointer added.
		// ie. initialized empty root node with no data.
		n.pointer = p
		return
	}

	// leaf node with a point in it.  Now we're splitting it and making it an internal node.
	nPoint := n.pointer.Point()
	n.internal = true

	if dx, dy := nPoint.X()-point.X(), nPoint.Y()-point.Y(); dx*dx+dy*dy <= q.Threshold*q.Threshold {
		// similar/duplicate point to stop recursion.
		i := childIndex((left+right)/2.0, (bottom+top)/2.0, point)
		if n.children[i] == nil {
			n.children[i] = q.nextNode()
			n.children[i].pointer = p
			return
		}
		q.insert(n, p, left, right, bottom, top)
		return
	}

	nPointer := n.pointer
	n.pointer = nil

	// current node is now an internal node.
	// first re add its pointer as one of its children,
	// then add the new pointer as one of the children.
	q.insert(n, nPointer, left, right, bottom, top)
	q.insert(n, p, left, right, bottom, top)
}

// Find returns the closest Value/Pointer in the quadtree.
// This function is thread safe. Multiple goroutines can read from
// a pre-created tree.
func (q *Quadtree) Find(p *geo.Point) geo.Pointer {
	if q.root == nil {
		return nil
	}

	v := &findVisitor{
		point:          p,
		closestBound:   q.bound.Clone(),
		minDistSquared: math.MaxFloat64,
	}

	newVisit(v).Visit(q.root,
		q.bound.Left(), q.bound.Right(),
		q.bound.Bottom(), q.bound.Top(),
	)

	return v.closest
}

// InBound returns a slice with all the pointers in the quadtree that are
// within the given bound. An optional buffer parameter is provided to allow
// for the reuse of result slice memory. This function is thread safe.
// Multiple goroutines can read from a pre-created tree.
func (q *Quadtree) InBound(b *geo.Bound, buf ...[]geo.Pointer) []geo.Pointer {
	if q.root == nil {
		return nil
	}

	var p []geo.Pointer
	if len(buf) > 0 {
		p = buf[0][:0]
	}
	v := &inBoundVisitor{
		bound:    b,
		pointers: p,
	}

	newVisit(v).Visit(q.root,
		q.bound.Left(), q.bound.Right(),
		q.bound.Bottom(), q.bound.Top(),
	)

	return v.pointers
}

// The visit stuff is a more go like (hopefully) implementation of the
// d3.quadtree.visit function. It is not exported, but if there is a
// good use case, it could be.

type visitor interface {
	// Bound returns the current relevant bound so we can prune irrelevant nodes
	// from the search.
	Bound() *geo.Bound
	Visit(p geo.Pointer)

	// Point should return the specific point being search for, or null if there
	// isn't one (ie. searching by bound). This helps guide the search to the
	// best child node first.
	Point() *geo.Point
}

// visit provides a framework for walking the quad tree.
// Currently used by the `Find` and `InBound` functions.
type visit struct {
	visitor visitor
}

func newVisit(v visitor) *visit {
	return &visit{
		visitor: v,
	}
}

func (v *visit) Visit(n *node, left, right, bottom, top float64) {
	b := v.visitor.Bound()
	if left > b.Right() || right < b.Left() ||
		bottom > b.Top() || top < b.Bottom() {
		return
	}

	if n.pointer != nil {
		v.visitor.Visit(n.pointer)
	}

	if !n.internal {
		return
	}

	cx := (left + right) / 2.0
	cy := (bottom + top) / 2.0

	i := 0
	if p := v.visitor.Point(); p != nil {
		// go to the child node the point is in first.
		i = childIndex(cx, cy, p)
	}

	for j := i; j < i+4; j++ {
		if n.children[j%4] == nil {
			continue
		}

		if k := j % 4; k == 0 {
			v.Visit(n.children[0], left, cx, cy, top)
		} else if k == 1 {
			v.Visit(n.children[1], cx, right, cy, top)
		} else if k == 2 {
			v.Visit(n.children[2], left, cx, bottom, cy)
		} else if k == 3 {
			v.Visit(n.children[3], cx, right, bottom, cy)
		}
	}
}

type findVisitor struct {
	point          *geo.Point
	closest        geo.Pointer
	closestBound   *geo.Bound
	minDistSquared float64
}

func (v *findVisitor) Bound() *geo.Bound {
	return v.closestBound
}

func (v *findVisitor) Point() *geo.Point {
	return v.point
}

func (v *findVisitor) Visit(p geo.Pointer) {
	point := p.Point()
	if d := point.SquaredDistanceFrom(v.point); d < v.minDistSquared {
		v.minDistSquared = d
		v.closest = p

		d = math.Sqrt(d)
		x := v.point.X()
		y := v.point.Y()
		v.closestBound.Set(x-d, x+d, y-d, y+d)
	}

	return
}

type inBoundVisitor struct {
	bound    *geo.Bound
	pointers []geo.Pointer
}

func (v *inBoundVisitor) Bound() *geo.Bound {
	return v.bound
}

func (v *inBoundVisitor) Point() *geo.Point {
	return nil
}

func (v *inBoundVisitor) Visit(p geo.Pointer) {
	if v.bound.Contains(p.Point()) {
		v.pointers = append(v.pointers, p)
	}
}

func childIndex(cx, cy float64, point *geo.Point) int {
	i := 0
	if point.Y() <= cy {
		i = 2
	}

	if point.X() >= cx {
		i++
	}

	return i
}
