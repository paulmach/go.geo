package quadtree

import (
	"math/rand"
	"testing"

	"github.com/paulmach/go.geo"
)

func TestNew(t *testing.T) {
	bound := geo.NewBound(0, 1, 2, 3)
	qt := New(bound)

	if qt.Bound() != bound {
		t.Errorf("should use provided bound, got %v", qt.Bound)
	}

	ps := geo.NewPointSet()
	ps.Push(geo.NewPoint(0, 2))
	ps.Push(geo.NewPoint(1, 3))

	qt = NewFromPointSet(ps)
	if !qt.Bound().Equals(ps.Bound()) {
		t.Errorf("should take bound from pointset, got %v", qt.Bound())
	}
}

func TestQuadtreeFind(t *testing.T) {
	points := []geo.Pointer{}
	dim := 17
	for i := 0; i < dim*dim; i++ {
		points = append(points, geo.NewPoint(float64(i%dim), float64(i/dim)))
	}

	q := NewFromPointers(points)

	// table test
	type findTest struct {
		Point    *geo.Point
		Expected *geo.Point
	}

	tests := []findTest{
		{Point: geo.NewPoint(0.1, 0.1), Expected: geo.NewPoint(0, 0)},
		{Point: geo.NewPoint(3.1, 2.9), Expected: geo.NewPoint(3, 3)},
		{Point: geo.NewPoint(7.5, 7.5), Expected: geo.NewPoint(7, 7)},
		{Point: geo.NewPoint(7.1, 7.1), Expected: geo.NewPoint(7, 7)},
		{Point: geo.NewPoint(8.5, 8.5), Expected: geo.NewPoint(7, 7)},
		{Point: geo.NewPoint(0.1, 15.9), Expected: geo.NewPoint(0, 16)},
		{Point: geo.NewPoint(15.9, 15.9), Expected: geo.NewPoint(16, 16)},
	}

	for i, test := range tests {
		if i != 3 {
			continue
		}
		if v := q.Find(test.Point); !v.Point().Equals(test.Expected) {
			t.Errorf("incorrect point on %d, got %v", i, v)
		}
	}
}

func TestQuadtreeFindRandom(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	ps := geo.NewPointSet()
	for i := 0; i < 1000; i++ {
		ps.Push(geo.NewPoint(r.Float64(), r.Float64()))
	}
	qt := NewFromPointSet(ps)

	for i := 0; i < 1000; i++ {
		p := geo.NewPoint(r.Float64(), r.Float64())

		f := qt.Find(p)
		_, j := ps.DistanceFrom(p)

		if e := ps.GetAt(j); !e.Equals(f.Point()) {
			t.Errorf("index: %d, unexpected point %v != %v", i, e, f.Point())
		}
	}
}

func TestQuadtreeInBoundRandom(t *testing.T) {
	r := rand.New(rand.NewSource(43))

	var pointers []geo.Pointer
	for i := 0; i < 1000; i++ {
		pointers = append(pointers, geo.NewPoint(r.Float64(), r.Float64()))
	}
	qt := NewFromPointers(pointers)

	for i := 0; i < 1000; i++ {
		p := geo.NewPoint(r.Float64(), r.Float64())

		b := geo.NewBoundFromPoints(p, p).Pad(0.1)
		ps := qt.InBound(b)

		// find the right answer brute force
		var list []geo.Pointer
		for _, p := range pointers {
			if b.Contains(p.Point()) {
				list = append(list, p)
			}
		}

		if len(list) != len(ps) {
			t.Errorf("index: %d, lengths not equal %v != %v", i, len(list), len(ps))
		}
	}
}
