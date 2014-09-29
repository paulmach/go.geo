package reducers

import (
	"testing"

	"github.com/paulmach/go.geo"
)

func TestVisvalingamThreshold(t *testing.T) {
	p := geo.NewPath()
	p.Push(geo.NewPoint(0.0, 0.0))
	p.Push(geo.NewPoint(1.0, 1.0))
	p.Push(geo.NewPoint(0.0, 2.0))
	p.Push(geo.NewPoint(1.0, 3.0))
	p.Push(geo.NewPoint(0.0, 4.0))

	reduced := VisvalingamThreshold(p, 1.1) // should reduce
	if l := reduced.Length(); l != 2 {
		t.Errorf("visvalingamThreshold reduce to incorrect number of points, expected 2, got %d", l)
	}

	reduced = VisvalingamThreshold(p, 0.9) // should not reduce
	if l := reduced.Length(); l != 5 {
		t.Errorf("visvalingamThreshold reduce to incorrect number of points, expected 5, got %d", l)
	}
}

func TestVisvalingamKeep(t *testing.T) {
	p := geo.NewPath()
	p.Push(geo.NewPoint(0.0, 0.0))
	p.Push(geo.NewPoint(1.0, 1.0))
	p.Push(geo.NewPoint(0.0, 2.0))
	p.Push(geo.NewPoint(1.0, 3.0))
	p.Push(geo.NewPoint(0.0, 4.0))

	for i := 6; i <= 7; i++ {
		reduced := VisvalingamKeep(p, i)
		if l := reduced.Length(); l != 5 {
			t.Errorf("visvalingamKeep reduce to incorrect number of points, expected %d, got %d", 5, l)
		}
	}

	for i := 2; i <= 5; i++ {
		reduced := VisvalingamKeep(p, i)
		if l := reduced.Length(); l != i {
			t.Errorf("visvalingamKeep reduce to incorrect number of points, expected %d, got %d", i, l)
		}
	}
}

func TestVisvalingam(t *testing.T) {
	p := geo.NewPath()
	reduced := Visvalingam(p, 0.1, 0)
	if !reduced.Equals(p) {
		t.Error("visvalingam should return same path if of length 0")
	}

	if reduced == p {
		t.Error("should create new path and not modify original")
	}

	p.Push(geo.NewPoint(0.0, 0.0))
	reduced = Visvalingam(p, 0.1, 0)
	if !reduced.Equals(p) {
		t.Error("visvalingam should return same path if of length 1")
	}

	if reduced == p {
		t.Error("should create new path and not modify original")
	}

	p.Push(geo.NewPoint(1.0, 1.0))
	reduced = Visvalingam(p, 0.1, 0)
	if !reduced.Equals(p) {
		t.Error("visvalingam should return same path if of length 2")
	}

	if reduced == p {
		t.Error("should create new path and not modify original")
	}

	// 3 points
	p.Push(geo.NewPoint(0.0, 2.0))

	reduced = Visvalingam(p, 1.1, 0) // should reduce
	if l := reduced.Length(); l != 2 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 2, got %d", l)
	}

	reduced = Visvalingam(p, 1.1, 3) // should not reduce
	if l := reduced.Length(); l != 3 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 3, got %d", l)
	}

	reduced = Visvalingam(p, 0.9, 0) // should not reduce
	if l := reduced.Length(); l != 3 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 3, got %d", l)
	}

	if reduced == p {
		t.Error("should create new path and not modify original")
	}

	// 5 points
	p.Push(geo.NewPoint(1.0, 3.0))
	p.Push(geo.NewPoint(0.0, 4.0))

	reduced = Visvalingam(p, 1.1, 0) // should reduce
	if l := reduced.Length(); l != 2 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 2, got %d", l)
	}

	reduced = Visvalingam(p, 1.1, 5) // should not reduce
	if l := reduced.Length(); l != 5 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 5, got %d", l)
	}

	reduced = Visvalingam(p, 1.1, 3) // should reduce
	if l := reduced.Length(); l != 3 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 3, got %d", l)
	}

	reduced = Visvalingam(p, 0.9, 0) // should not reduce
	if l := reduced.Length(); l != 5 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 5, got %d", l)
	}

	if reduced == p {
		t.Error("should create new path and not modify original")
	}

	// colinear points
	p = geo.NewPath()
	p.Push(geo.NewPoint(0, 0))
	p.Push(geo.NewPoint(0, 1))
	p.Push(geo.NewPoint(0, 2))

	if l := Visvalingam(p, 0.0, 0).Length(); l != 2 {
		t.Errorf("visvalingam reduce should remove coliniar points")
	}
}

func TestVisvalingamPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("path, intersection invalid geometry should panic")
		}
	}()

	// this should panic
	Visvalingam(geo.NewPath(), -100, 0)
}

func TestDoubleTriangleArea(t *testing.T) {
	p1 := geo.NewPoint(2, 5)
	p2 := geo.NewPoint(5, 1)
	p3 := geo.NewPoint(-4, 3)

	expected := 30.0

	// check all the orderings
	if area := doubleTriangleArea(p1, p2, p3); area != expected {
		t.Errorf("triangleArea expected %f, got %f", expected, area)
	}

	if area := doubleTriangleArea(p1, p3, p2); area != expected {
		t.Errorf("triangleArea expected %f, got %f", expected, area)
	}

	if area := doubleTriangleArea(p2, p3, p1); area != expected {
		t.Errorf("triangleArea expected %f, got %f", expected, area)
	}

	if area := doubleTriangleArea(p2, p1, p3); area != expected {
		t.Errorf("triangleArea expected %f, got %f", expected, area)
	}

	if area := doubleTriangleArea(p3, p1, p2); area != expected {
		t.Errorf("triangleArea expected %f, got %f", expected, area)
	}

	if area := doubleTriangleArea(p3, p2, p1); area != expected {
		t.Errorf("triangleArea expected %f, got %f", expected, area)
	}
}
