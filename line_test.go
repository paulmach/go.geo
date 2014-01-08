package geo

import (
	"testing"
)

func TestLineNew(t *testing.T) {
	a := NewPoint(1, 1)
	b := NewPoint(2, 2)

	l := NewLine(a, b)
	if !l.A().Equals(a) {
		t.Errorf("line, expected %v == %v", l.A(), a)
	}

	if !l.B().Equals(b) {
		t.Errorf("line, expected %v == %v", l.B(), b)
	}

	// verify there is a clone
	b.Scale(10)
	if l.B().Equals(b) {
		t.Errorf("line, expected %v != %v", l.B(), b)
	}
}

func TestLineDistanceFrom(t *testing.T) {
	var answer float64
	l := NewLine(NewPoint(0, 0), NewPoint(0, 10))

	answer = 1
	if d := l.DistanceFrom(NewPoint(1, 5)); d != answer {
		t.Errorf("line, distanceFrom expected %f, got %f", answer, d)
	}

	answer = 0
	if d := l.DistanceFrom(NewPoint(0, 2)); d != answer {
		t.Errorf("line, distanceFrom expected %f, got %f", answer, d)
	}

	answer = 5
	if d := l.DistanceFrom(NewPoint(0, -5)); d != answer {
		t.Errorf("line, distanceFrom expected %f, got %f", answer, d)
	}

	answer = 3
	if d := l.DistanceFrom(NewPoint(0, 13)); d != answer {
		t.Errorf("line, distanceFrom expected %f, got %f", answer, d)
	}

	l = NewLine(NewPoint(0, 0), NewPoint(0, 0))
	answer = 5
	if d := l.DistanceFrom(NewPoint(3, 4)); d != answer {
		t.Errorf("line, distanceFrom expected %f, got %f", answer, d)
	}
}

func TestLineDistance(t *testing.T) {
	l := NewLine(NewPoint(0, 0), NewPoint(3, 4))
	if d := l.Distance(); d != 5 {
		t.Errorf("line, distance expected 5, got %f", d)
	}

	l.B().Scale(2)
	if d := l.Distance(); d != 10 {
		t.Errorf("line, distance expected 10, got %f", d)
	}
}

func TestLineInterpolate(t *testing.T) {
	var answer *Point
	l := NewLine(NewPoint(0, 0), NewPoint(5, 10))

	answer = NewPoint(1, 2)
	if p := l.Interpolate(.20); !p.Equals(answer) {
		t.Errorf("line, interpolate expected %v, got %v", answer, p)
	}

	answer = NewPoint(4, 8)
	if p := l.Interpolate(.80); !p.Equals(answer) {
		t.Errorf("line, interpolate expected %v, got %v", answer, p)
	}

	answer = NewPoint(-1, -2)
	if p := l.Interpolate(-.20); !p.Equals(answer) {
		t.Errorf("line, interpolate expected %v, got %v", answer, p)
	}

	answer = NewPoint(6, 12)
	if p := l.Interpolate(1.20); !p.Equals(answer) {
		t.Errorf("line, interpolate expected %v, got %v", answer, p)
	}
}

func TestLineSide(t *testing.T) {
	l := NewLine(NewPoint(0, 0), NewPoint(0, 10))

	// colinear
	if o := l.Side(NewPoint(0, -5)); o != 0 {
		t.Errorf("point, expected to be colinear, got %d", o)
	}

	// right
	if o := l.Side(NewPoint(1, 5)); o != 1 {
		t.Errorf("point, expected to be on right, got %d", o)
	}

	// left
	if o := l.Side(NewPoint(-1, 5)); o != -1 {
		t.Errorf("point, expected to be on left, got %d", o)
	}
}

func TestLineIntersection(t *testing.T) {
	var answer *Point
	l := NewLine(NewPoint(0, 0), NewPoint(1, 1))

	answer = nil
	if p := l.Intersection(NewLine(NewPoint(1, 0), NewPoint(2, 1))); p != nil {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = nil
	if p := l.Intersection(NewLine(NewPoint(1, 0), NewPoint(3, 1))); p != nil {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = InfinityPoint
	if p := l.Intersection(NewLine(NewPoint(1, 1), NewPoint(2, 2))); !p.Equals(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = NewPoint(1, 1)
	if p := l.Intersection(NewLine(NewPoint(1, 1), NewPoint(2, 3))); !p.Equals(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = NewPoint(0, 0)
	if p := l.Intersection(NewLine(NewPoint(1, 10), NewPoint(0, 0))); !p.Equals(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = NewPoint(0.5, 0.5)
	if p := l.Intersection(NewLine(NewPoint(0, 1), NewPoint(1, 0))); !p.Equals(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = NewPoint(0.5, 0.5)
	if p := l.Intersection(NewLine(NewPoint(0, 1), NewPoint(2, -1))); !p.Equals(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}

	answer = NewPoint(0.5, 0.5)
	if p := l.Intersection(NewLine(NewPoint(0.5, 0.5), NewPoint(2, -1))); !p.Equals(answer) {
		t.Errorf("line, intersection expected %v, got %v", answer, p)
	}
}

func TestLineIntersects(t *testing.T) {
	var answer bool
	l := NewLine(NewPoint(0, 0), NewPoint(1, 1))

	answer = false
	if p := l.Intersects(NewLine(NewPoint(1, 0), NewPoint(2, 1))); p != answer {
		t.Errorf("line, intersects expected %v, got %v", answer, p)
	}

	answer = true
	if p := l.Intersects(NewLine(NewPoint(1, 0), NewPoint(0, 1))); p != answer {
		t.Errorf("line, intersects expected %v, got %v", answer, p)
	}

	answer = true
	if p := l.Intersects(NewLine(NewPoint(1, 1), NewPoint(2, 1))); p != answer {
		t.Errorf("line, intersects expected %v, got %v", answer, p)
	}

	answer = true
	l2 := NewLine(NewPoint(0.5, 0.5), NewPoint(2, 2))
	if p := l.Intersects(l2); p != answer {
		t.Errorf("line, intersects expected %v, got %v", answer, p)
	}

	if p := l2.Intersects(l); p != answer {
		t.Errorf("line, intersects expected %v, got %v", answer, p)
	}
}

func TestLineMidpoint(t *testing.T) {
	var answer *Point
	l := NewLine(NewPoint(0, 0), NewPoint(10, 20))

	answer = NewPoint(5, 10)
	if p := l.Midpoint(); !p.Equals(answer) {
		t.Errorf("line, midpoint expected %v, got %v", answer, p)
	}
}

func TestLineBounds(t *testing.T) {
	var answer *Bound
	a := NewPoint(1, 2)
	b := NewPoint(3, 4)

	l := NewLine(a, b)

	answer = NewBound(1, 3, 2, 4)
	if b := l.Bounds(); !b.Equals(answer) {
		t.Errorf("line, bounds expected %v, got %v", answer, b)
	}

	if b := l.Reverse().Bounds(); !b.Equals(answer) {
		t.Errorf("line, bounds expected %v, got %v", answer, b)
	}
}

func TestLineReverse(t *testing.T) {
	a := NewPoint(1, 2)
	b := NewPoint(3, 4)

	l := NewLine(a, b).Reverse()

	if !l.A().Equals(b) || !l.B().Equals(a) {
		t.Error("line, reverse did not work")
	}
}

func TestLineClone(t *testing.T) {
	l1 := NewLine(NewPoint(1, 1), NewPoint(2, 2))
	l2 := l1.Clone()

	l1.A().Scale(10)
	l2.B().Scale(15)

	if l1.A().Equals(l2.A()) {
		t.Errorf("line, clone expected %v != %v", l1.A(), l2.A())
	}

	if l2.B().Equals(l1.B()) {
		t.Errorf("line, clone expected %v != %v", l2.B(), l1.B())
	}
}

func TestLineEquals(t *testing.T) {
	l1 := NewLine(NewPoint(1, 2), NewPoint(3, 4))
	l2 := NewLine(NewPoint(1, 2), NewPoint(3, 4))

	if !l1.Equals(l2) || !l2.Equals(l1) {
		t.Errorf("line, equals expcted %v == %v", l1, l2)
	}

	l3 := NewLine(NewPoint(3, 4), NewPoint(1, 2))
	if !l1.Equals(l3) || !l3.Equals(l1) {
		t.Errorf("line, equals expcted %v == %v", l1, l3)
	}
}
