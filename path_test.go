package geo

import (
	"math"
	"math/rand"
	"testing"
)

func TestPathReduce(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(0.5, .2))
	p.Push(NewPoint(1, 0))

	if l := p.Clone().Reduce(0.1).Length(); l != 3 {
		t.Errorf("path, reduce to incorrect number of points, expected 2, got %d", l)
	}

	if l := p.Clone().Reduce(0.3).Length(); l != 2 {
		t.Errorf("path, reduce to incorrect number of points, expected 3, got %d", l)
	}
}

func TestPathEncode(t *testing.T) {
	for loop := 0; loop < 100; loop++ {
		p := NewPath()
		for i := 0; i < 100; i++ {
			p.Push(&Point{rand.Float64(), rand.Float64()})
		}

		encoded := p.Encode()
		for _, c := range encoded {
			if c < 63 || c > 127 {
				t.Errorf("path, encode result out of range: %d", c)
			}
		}
	}
}

func TestPathEncodeDecode(t *testing.T) {
	for loop := 0; loop < 100; loop++ {

		p := NewPath()
		for i := 0; i < 100; i++ {
			p.Push(&Point{rand.Float64(), rand.Float64()})
		}

		encoded := p.Encode(int(1.0 / epsilon))
		path := Decode(encoded, int(1.0/epsilon))

		if path.Length() != 100 {
			t.Fatalf("path, encodeDecode length mismatch: %d != 100", path.Length())
		}

		for i := 0; i < 100; i++ {
			a := p.GetAt(i)
			b := path.GetAt(i)

			if e := math.Abs(a[0] - b[0]); e > epsilon {
				t.Errorf("path, encodeDecode X error too big: %f", e)
			}

			if e := math.Abs(a[1] - b[1]); e > epsilon {
				t.Errorf("path, encodeDecode Y error too big: %f", e)
			}
		}
	}
}

func TestPathDistance(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(0, 3))
	p.Push(NewPoint(4, 3))

	if d := p.Distance(); d != 7 {
		t.Errorf("path, distance got: %f, expected 7.0", d)
	}
}

func TestPathDistanceFrom(t *testing.T) {
	var answer float64

	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(0, 3))
	p.Push(NewPoint(4, 3))
	p.Push(NewPoint(4, 0))

	answer = 0.5
	if d := p.DistanceFrom(NewPoint(4.5, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, distanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.4
	if d := p.DistanceFrom(NewPoint(0.4, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, distanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.3
	if d := p.DistanceFrom(NewPoint(-0.3, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, distanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.2
	if d := p.DistanceFrom(NewPoint(0.3, 2.8)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, distanceFrom expected %f, got: %f", answer, d)
	}
}

func TestPathBound(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0.5, .2))
	p.Push(NewPoint(-1, 0))
	p.Push(NewPoint(1, 10))
	p.Push(NewPoint(1, 8))

	answer := NewBound(-1, 1, 0, 10)
	if b := p.Bound(); !b.Equals(answer) {
		t.Errorf("path, bound, %v != %v", b, answer)
	}
}

func TestPathEquals(t *testing.T) {
	p1 := NewPath()
	p1.Push(NewPoint(0.5, .2))
	p1.Push(NewPoint(-1, 0))
	p1.Push(NewPoint(1, 10))

	p2 := NewPath()
	p2.Push(NewPoint(0.5, .2))
	p2.Push(NewPoint(-1, 0))
	p2.Push(NewPoint(1, 10))

	if !p1.Equals(p2) {
		t.Errorf("path, equals paths should be equal")
	}

	p3 := p2.Clone().SetAt(1, NewPoint(0, 0))
	if p1.Equals(p3) {
		t.Errorf("path, equals paths should not be equal")
	}

	p2.Pop()
	if p2.Equals(p1) {
		t.Errorf("path, equals paths should not be equal")
	}
}

func TestPathClone(t *testing.T) {
	p1 := NewPath()
	p1.Push(NewPoint(0, 0))
	p1.Push(NewPoint(0.5, .2))
	p1.Push(NewPoint(1, 0))

	p2 := p1.Clone()
	p2.Pop()
	if p1.Length() == p2.Length() {
		t.Errorf("path, clone length %d == %d", p1.Length(), p2.Length())
	}

	p2 = p1.Clone()
	if p1 == p2 {
		t.Errorf("path, clone should return different pointers")
	}

	if !p2.Equals(p1) {
		t.Errorf("path, clone paths should be equal")
	}
}
