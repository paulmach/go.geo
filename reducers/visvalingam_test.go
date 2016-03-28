package reducers

import (
	"testing"

	"github.com/paulmach/go.geo"
)

func TestVisvalingamThreshold(t *testing.T) {
	p := append(geo.NewPath(),
		geo.NewPoint(0.0, 0.0),
		geo.NewPoint(1.0, 1.0),
		geo.NewPoint(0.0, 2.0),
		geo.NewPoint(1.0, 3.0),
		geo.NewPoint(0.0, 4.0),
	)

	reduced := VisvalingamThreshold(p, 1.1) // should reduce
	if l := len(reduced); l != 2 {
		t.Errorf("visvalingamThreshold reduce to incorrect number of points, expected 2, got %d", l)
	}

	reduced = VisvalingamThreshold(p, 0.9) // should not reduce
	if l := len(reduced); l != 5 {
		t.Errorf("visvalingamThreshold reduce to incorrect number of points, expected 5, got %d", l)
	}
}

func TestVisvalingamKeep(t *testing.T) {
	p := append(geo.NewPath(),
		geo.NewPoint(0.0, 0.0),
		geo.NewPoint(1.0, 1.0),
		geo.NewPoint(0.0, 2.0),
		geo.NewPoint(1.0, 3.0),
		geo.NewPoint(0.0, 4.0),
	)

	for i := 6; i <= 7; i++ {
		reduced := VisvalingamKeep(p, i)
		if l := len(reduced); l != 5 {
			t.Errorf("visvalingamKeep reduce to incorrect number of points, expected %d, got %d", 5, l)
		}
	}

	for i := 2; i <= 5; i++ {
		reduced := VisvalingamKeep(p, i)
		if l := len(reduced); l != i {
			t.Errorf("visvalingamKeep reduce to incorrect number of points, expected %d, got %d", i, l)
		}
	}
}

func TestVisvalingam(t *testing.T) {
	p := geo.NewPath()
	reduced := Visvalingam(p, 0.1, 0)
	if !reduced.Equal(p) {
		t.Error("visvalingam should return same path if of length 0")
	}

	p = append(p, geo.NewPoint(0.0, 0.0))
	reduced = Visvalingam(p, 0.1, 0)
	if !reduced.Equal(p) {
		t.Error("visvalingam should return same path if of length 1")
	}

	p = append(p, geo.NewPoint(1.0, 1.0))
	reduced = Visvalingam(p, 0.1, 0)
	if !reduced.Equal(p) {
		t.Error("visvalingam should return same path if of length 2")
	}

	// 3 points
	p = append(p, geo.NewPoint(0.0, 2.0))
	reduced = Visvalingam(p, 1.1, 0) // should reduce
	if l := len(reduced); l != 2 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 2, got %d", l)
	}

	reduced = Visvalingam(p, 1.1, 3) // should not reduce
	if l := len(reduced); l != 3 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 3, got %d", l)
	}

	reduced = Visvalingam(p, 0.9, 0) // should not reduce
	if l := len(reduced); l != 3 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 3, got %d", l)
	}

	// 5 points
	p = append(p,
		geo.NewPoint(1.0, 3.0),
		geo.NewPoint(0.0, 4.0),
	)

	reduced = Visvalingam(p, 1.1, 0) // should reduce
	if l := len(reduced); l != 2 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 2, got %d", l)
	}

	reduced = Visvalingam(p, 1.1, 5) // should not reduce
	if l := len(reduced); l != 5 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 5, got %d", l)
	}

	reduced = Visvalingam(p, 1.1, 3) // should reduce
	if l := len(reduced); l != 3 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 3, got %d", l)
	}

	reduced = Visvalingam(p, 0.9, 0) // should not reduce
	if l := len(reduced); l != 5 {
		t.Errorf("visvalingam reduce to incorrect number of points, expected 5, got %d", l)
	}

	// colinear points
	p = append(geo.NewPath(),
		geo.NewPoint(0, 0),
		geo.NewPoint(0, 1),
		geo.NewPoint(0, 2),
	)

	if l := len(Visvalingam(p, 0.0, 0)); l != 2 {
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
	pp1 := geo.NewPoint(2, 5)
	pp2 := geo.NewPoint(5, 1)
	pp3 := geo.NewPoint(-4, 3)
	p1 := &pp1
	p2 := &pp2
	p3 := &pp3

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
