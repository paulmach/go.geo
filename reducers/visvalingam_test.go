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

func TestSqDistance(t *testing.T) {
	expected := 25.0
	if d := sqDistance(geo.NewPoint(0, 3), geo.NewPoint(4, 0)); d != expected {
		t.Errorf("sqDistance expected %f, got %f", expected, d)
	}
}

func TestTrianglePointNormalArea(t *testing.T) {
	expected := 576.0
	if d := trianglePointNormalArea(geo.NewPoint(0, 3), geo.NewPoint(4, 0), geo.NewPoint(0, 0)); d != expected {
		t.Errorf("trianglePointNormalArea expected %f, got %f", expected, d)
	}
}

func TestTriangleSquareDistanceNormalArea(t *testing.T) {
	expected := 576.0
	if d := triangleSquareDistanceNormalArea(9, 16, 25); d != expected {
		t.Errorf("triangleSquareDistanceNormalArea expected %f, got %f", expected, d)
	}
}
