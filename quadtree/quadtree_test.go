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
		t.Errorf("should use provided bound, got %v", qt.Bound())
	}

	if qt.freeNodes != nil {
		t.Errorf("freeNodes should not be preallocated")
	}

	qt = New(bound, 12)

	if len(qt.freeNodes) != 12 {
		t.Errorf("should preallocate %d freeNodes", 12)
	}

	ps := geo.NewPointSet()
	ps.Push(geo.NewPoint(0, 2))
	ps.Push(geo.NewPoint(1, 3))

	qt = NewFromPointSet(ps)
	if !qt.Bound().Equals(ps.Bound()) {
		t.Errorf("should take bound from pointset, got %v", qt.Bound())
	}

	if len(qt.freeNodes) != ps.Length() {
		t.Errorf("should preallocate %d freeNodes", ps.Length())
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

func TestQuadtreeFindMatching(t *testing.T) {
	type dataPointer struct {
		geo.Pointer
		visible bool
	}

	q := NewFromPointers([]geo.Pointer{
		dataPointer{geo.NewPoint(0, 0), false},
		dataPointer{geo.NewPoint(1, 1), true},
	})

	// filters
	filters := map[bool]Filter{
		false: nil,
		true:  func(p geo.Pointer) bool { return p.(dataPointer).visible },
	}

	// table test
	type findTest struct {
		Filtered bool
		Point    *geo.Point
		Expected *geo.Point
	}

	tests := []findTest{
		{Filtered: false, Point: geo.NewPoint(0.1, 0.1), Expected: geo.NewPoint(0, 0)},
		{Filtered: true, Point: geo.NewPoint(0.1, 0.1), Expected: geo.NewPoint(1, 1)},
	}

	for i, test := range tests {
		if v := q.FindMatching(test.Point, filters[test.Filtered]); !v.Point().Equals(test.Expected) {
			t.Errorf("incorrect point on %d, got %v", i, v)
		}
	}
}

func TestQuadtreeFindKNearest(t *testing.T) {
	type dataPointer struct {
		geo.Pointer
		visible bool
	}

	q := NewFromPointers([]geo.Pointer{
		dataPointer{geo.NewPoint(0, 0), false},
		dataPointer{geo.NewPoint(1, 1), true},
		dataPointer{geo.NewPoint(2, 2), false},
		dataPointer{geo.NewPoint(3, 3), true},
		dataPointer{geo.NewPoint(4, 4), false},
		dataPointer{geo.NewPoint(5, 5), true},
	})

	// filters
	filters := map[bool]Filter{
		false: nil,
		true:  func(p geo.Pointer) bool { return p.(dataPointer).visible },
	}

	// table test
	type findTest struct {
		Filtered bool
		Point    *geo.Point
		Expected []*geo.Point
	}

	tests := []findTest{
		{Filtered: false, Point: geo.NewPoint(0.1, 0.1), Expected: []*geo.Point{geo.NewPoint(0, 0), geo.NewPoint(1, 1)}},
		{Filtered: true, Point: geo.NewPoint(0.1, 0.1), Expected: []*geo.Point{geo.NewPoint(1, 1), geo.NewPoint(3, 3)}},
	}

	for i, test := range tests {
		v := q.FindKNearestMatching(test.Point, 2, filters[test.Filtered])
		if len(v) != len(test.Expected) {
			t.Errorf("incorrect response length on %d", i)
		}
		for _, answer := range v {
			found := false
			for _, expected := range test.Expected {
				if answer.Point().Equals(expected) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("incorrect point on %d, got %v", i, v)
			}
		}
	}
}

func TestQuadtreeFindKNearestWithDistanceLimit(t *testing.T) {
	type dataPointer struct {
		geo.Pointer
		visible bool
	}

	q := NewFromPointers([]geo.Pointer{
		dataPointer{geo.NewPoint(0, 0), false},
		dataPointer{geo.NewPoint(1, 1), true},
		dataPointer{geo.NewPoint(2, 2), false},
		dataPointer{geo.NewPoint(3, 3), true},
		dataPointer{geo.NewPoint(4, 4), false},
		dataPointer{geo.NewPoint(5, 5), true},
	})

	// filters
	filters := map[bool]Filter{
		false: nil,
		true:  func(p geo.Pointer) bool { return p.(dataPointer).visible },
	}

	// table test
	type findTest struct {
		Filtered bool
		Distance float64
		Point    *geo.Point
		Expected []*geo.Point
	}

	tests := []findTest{
		{Filtered: false, Distance: 2., Point: geo.NewPoint(0.1, 0.1), Expected: []*geo.Point{geo.NewPoint(0, 0), geo.NewPoint(1, 1)}},
		{Filtered: true, Distance: 5., Point: geo.NewPoint(0.1, 0.1), Expected: []*geo.Point{geo.NewPoint(1, 1), geo.NewPoint(3, 3)}},
	}

	for i, test := range tests {
		v := q.FindKNearestMatching(test.Point, 10, filters[test.Filtered], test.Distance)
		if len(v) != len(test.Expected) {
			t.Errorf("incorrect response length on %d", i)
		}
		for _, answer := range v {
			found := false
			for _, expected := range test.Expected {
				if answer.Point().Equals(expected) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("incorrect point on %d, got %v", i, v)
			}
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

func TestQuadtreeInBoundMatching(t *testing.T) {
	type dataPointer struct {
		geo.Pointer
		visible bool
	}

	q := NewFromPointers([]geo.Pointer{
		dataPointer{geo.NewPoint(0, 0), false},
		dataPointer{geo.NewPoint(1, 1), true},
		dataPointer{geo.NewPoint(2, 2), false},
		dataPointer{geo.NewPoint(3, 3), true},
		dataPointer{geo.NewPoint(4, 4), false},
		dataPointer{geo.NewPoint(5, 5), true},
	})

	// filters
	filters := map[bool]Filter{
		false: nil,
		true:  func(p geo.Pointer) bool { return p.(dataPointer).visible },
	}

	// table test
	type findTest struct {
		Filtered bool
		Bound    *geo.Bound
		Expected []*geo.Point
	}

	tests := []findTest{
		{
			Filtered: false,
			Bound:    geo.NewBound(1, 3, 1, 3),
			Expected: []*geo.Point{
				geo.NewPoint(1, 1),
				geo.NewPoint(2, 2),
				geo.NewPoint(3, 3),
			},
		},
		{
			Filtered: true,
			Bound:    geo.NewBound(1, 3, 1, 3),
			Expected: []*geo.Point{
				geo.NewPoint(1, 1),
				geo.NewPoint(3, 3),
			},
		},
	}

	for i, test := range tests {
		v := q.InBoundMatching(test.Bound, filters[test.Filtered])
		if len(v) != len(test.Expected) {
			t.Errorf("incorrect response length on %d", i)
		}
		for _, answer := range v {
			found := false
			for _, expected := range test.Expected {
				if answer.Point().Equals(expected) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("incorrect point on %d, got %v", i, v)
			}
		}
	}
}
