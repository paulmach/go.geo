package geo

import (
	"testing"
)

func TestBoundNew(t *testing.T) {
	bound := NewBound(5, 0, 3, 0)
	if !bound.sw.Equals(&Point{0, 0}) {
		t.Errorf("bound, incorrect sw: expected %v, got %v", &Point{0, 0}, bound.sw)
	}

	if !bound.ne.Equals(&Point{5, 3}) {
		t.Errorf("bound, incorrect ne: expected %v, got %v", &Point{5, 3}, bound.ne)
	}

	bound = NewBoundFromPoints(&Point{0, 3}, &Point{4, 0})
	if !bound.sw.Equals(&Point{0, 0}) {
		t.Errorf("bound, incorrect sw: expected %v, got %v", &Point{0, 0}, bound.sw)
	}

	if !bound.ne.Equals(&Point{4, 3}) {
		t.Errorf("bound, incorrect ne: expected %v, got %v", &Point{4, 3}, bound.ne)
	}
}

func TestBoundExtend(t *testing.T) {
	bound := NewBound(3, 0, 5, 0)

	if b := bound.Clone().Extend(&Point{2, 1}); !b.Equals(bound) {
		t.Errorf("bound, extend: expected %v, got %v", bound, b)
	}

	answer := NewBound(6, 0, 5, -1)
	if b := bound.Clone().Extend(&Point{6, -1}); !b.Equals(answer) {
		t.Errorf("bound, extend: expected %v, got %v", answer, b)
	}
}

func TestBoundContains(t *testing.T) {
	var p *Point
	bound := NewBound(2, -2, 1, -1)

	p = &Point{0, 0}
	if !bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to be within %v", p, bound)
	}

	p = &Point{-1, 0}
	if !bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to be within %v", p, bound)
	}

	p = &Point{2, 1}
	if !bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to be within %v", p, bound)
	}

	p = &Point{0, 3}
	if bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to not be within %v", p, bound)
	}

	p = &Point{0, -3}
	if bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to not be within %v", p, bound)
	}

	p = &Point{3, 0}
	if bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to not be within %v", p, bound)
	}

	p = &Point{-3, 0}
	if bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to not be within %v", p, bound)
	}
}

func TestBoundIntersects(t *testing.T) {
	var tester *Bound
	bound := NewBound(0, 1, 2, 3)

	tester = NewBound(5, 6, 7, 8)
	if bound.Intersects(tester) {
		t.Errorf("bound, intersects expected %v, to not intersect %v", tester, bound)
	}

	tester = NewBound(-6, -5, 7, 8)
	if bound.Intersects(tester) {
		t.Errorf("bound, intersects expected %v, to not intersect %v", tester, bound)
	}

	tester = NewBound(0, 0.5, 7, 8)
	if bound.Intersects(tester) {
		t.Errorf("bound, intersects expected %v, to not intersect %v", tester, bound)
	}

	tester = NewBound(0, 0.5, 1, 4)
	if !bound.Intersects(tester) {
		t.Errorf("bound, intersects expected %v, to intersect %v", tester, bound)
	}

	tester = NewBound(-1, 2, 1, 4)
	if !bound.Intersects(tester) {
		t.Errorf("bound, intersects expected %v, to intersect %v", tester, bound)
	}

	tester = NewBound(0.3, 0.6, 2.3, 2.6)
	if !bound.Intersects(tester) {
		t.Errorf("bound, intersects expected %v, to intersect %v", tester, bound)
	}
}

func TestBoundCenter(t *testing.T) {
	var p *Point
	var b *Bound

	b = NewBound(0, 1, 2, 3)
	p = &Point{0.5, 2.5}
	if c := b.Center(); !c.Equals(p) {
		t.Errorf("bound, center expected %v, got %v", p, c)
	}

	b = NewBound(0, 0, 2, 2)
	p = &Point{0, 2}
	if c := b.Center(); !c.Equals(p) {
		t.Errorf("bound, center expected %v, got %v", p, c)
	}
}

func TestBoundPad(t *testing.T) {
	var bound, tester *Bound

	bound = NewBound(0, 1, 2, 3)
	tester = NewBound(-0.5, 1.5, 1.5, 3.5)
	if bound.Pad(0.5); !bound.Equals(tester) {
		t.Errorf("bound, pad expected %v, got %v", tester, bound)
	}

	bound = NewBound(0, 1, 2, 3)
	tester = NewBound(0.1, 0.9, 2.1, 2.9)
	if bound.Pad(-0.1); !bound.Equals(tester) {
		t.Errorf("bound, pad expected %v, got %v", tester, bound)
	}
}
