package quadtree

import (
	"github.com/paulmach/go.geo"
)

type Quadtree struct {
	TreeExtent *Extent
}

type QuadTreeNode struct {
	Leaf  bool
	Nodes []*QuadTreeNode
	Point *geo.Point
}

type Extent struct {
	TopLeft     *geo.Point
	BottomRight *geo.Point
}

type visitFunction func(*geo.Point, float64, float64, float64, float64) bool

func NewFromPoints(points *geo.PointSet) *Quadtree {
	return &Quadtree{}
}

func New() *Quadtree {
	return &Quadtree{}
}

func (q *Quadtree) Add(p *geo.Point) *Quadtree {
	return nil
}

func (q *Quadtree) Visit(visit visitFunction) {
	return nil
}

func (q *Quadtree) Find(p *geo.Point) *geo.Point {
	return nil
}

func (q *Quadtree) SetExtent(e *Extent) {
	q.TreeExtent = Extent
}
