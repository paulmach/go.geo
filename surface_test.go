package geo

import (
	"math"
	"testing"
)

type testValue struct {
	X, Y int
	A    *Point
}

func TestPointAt(t *testing.T) {
	bound := NewBound(3, 0, 3, 0)
	surface := NewSurface(bound, 7, 7)

	tests := []testValue{
		{0, 0, &Point{0, 0}}, {1, 1, &Point{0.5, 0.5}}, {2, 2, &Point{1, 1}}, {3, 3, &Point{1.5, 1.5}},
		{4, 5, &Point{2, 2.5}}, {5, 4, &Point{2.5, 2}}, {6, 0, &Point{3, 0}},
	}

	for _, point := range tests {
		if v := surface.PointAt(point.X, point.Y); *v != *point.A {
			t.Errorf("incorrect value at: expected %v, got %v", *point.A, *v)
		}
	}
}

type testPoint struct {
	P *Point
	A float64
}

func TestValueAt(t *testing.T) {
	bound := NewBound(3, 0, 3, 0)
	surface := NewSurface(bound, 4, 4)

	surface.Grid[1][1] = 0
	surface.Grid[2][1] = 1
	surface.Grid[1][2] = 2
	surface.Grid[2][2] = 3

	tests := []testPoint{
		{&Point{1, 1}, 0}, {&Point{2, 1}, 1}, {&Point{1, 2}, 2}, {&Point{2, 2}, 3},
		{&Point{1, 1.5}, 1}, {&Point{1.5, 1}, 0.5}, {&Point{2, 1.5}, 2}, {&Point{1.5, 2}, 2.5},
		{&Point{0, 0}, 0}, {&Point{1.5, 1.5}, 1.5}, {&Point{3, 3}, 0}, {&Point{6, 6}, 0},
		{&Point{1, 1.1}, 0.2}, {&Point{2, 1.1}, 1.2}, {&Point{1.1, 1.1}, 0.3},
	}

	for i, point := range tests {
		if v := surface.ValueAt(point.P); math.Abs(v-point.A) > epsilon {
			t.Errorf("surface: (%d) incorrect value at: expected %v, got %v", i, point.A, v)
		}
	}
}
