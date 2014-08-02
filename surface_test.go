package geo

import (
	"bytes"
	"math"
	"testing"
)

type testValue struct {
	X, Y int
	A    *Point
}

func TestSurfacePointAt(t *testing.T) {
	surface := NewSurface(NewBound(3, 0, 3, 0), 7, 7)

	tests := []testValue{
		{0, 0, &Point{0, 0}}, {1, 1, &Point{0.5, 0.5}}, {2, 2, &Point{1, 1}}, {3, 3, &Point{1.5, 1.5}},
		{4, 5, &Point{2, 2.5}}, {5, 4, &Point{2.5, 2}}, {6, 0, &Point{3, 0}},
	}

	for _, point := range tests {
		if point.A == nil {
			if v := surface.PointAt(point.X, point.Y); v != nil {
				t.Errorf("surface, pointAt incorrect value at: expected %v, got %v", nil, v)
			}
		} else {
			if v := surface.PointAt(point.X, point.Y); *v != *point.A {
				t.Errorf("surface, pointAt incorrect value at: expected %v, got %v", *point.A, *v)
			}
		}
	}
}

func TestSurfacePointAtPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("surface, pointAt out of range should panic")
		}
	}()

	surface := NewSurface(NewBound(3, 0, 3, 0), 7, 7)
	surface.PointAt(10, 10)
}

type testPointFloat struct {
	P *Point
	A float64
}

func TestSurfaceValueAt(t *testing.T) {
	bound := NewBound(3, 0, 3, 0)
	surface := NewSurface(bound, 4, 4)

	surface.Grid[1][1] = 0
	surface.Grid[2][1] = 1
	surface.Grid[1][2] = 2
	surface.Grid[2][2] = 3

	tests := []testPointFloat{
		{&Point{1, 1}, 0}, {&Point{2, 1}, 1}, {&Point{1, 2}, 2}, {&Point{2, 2}, 3},
		{&Point{1, 1.5}, 1}, {&Point{1.5, 1}, 0.5}, {&Point{2, 1.5}, 2}, {&Point{1.5, 2}, 2.5},
		{&Point{0, 0}, 0}, {&Point{1.5, 1.5}, 1.5}, {&Point{3, 3}, 0}, {&Point{6, 6}, 0},
		{&Point{1, 1.1}, 0.2}, {&Point{2, 1.1}, 1.2}, {&Point{1.1, 1.1}, 0.3},
	}

	// octave code
	// grid = [0, 0, 0, 0; 0, 0, 2, 0; 0, 1, 3, 0; 0, 0, 0, 0]'
	// interp2(0:3, 0:3, grid, 2, 1)

	for i, point := range tests {
		if v := surface.ValueAt(point.P); math.Abs(v-point.A) > epsilon {
			t.Errorf("surface, (%d) incorrect value at: expected %v, got %v", i, point.A, v)
		}
	}
}

type testPointPoint struct {
	P *Point
	A *Point
}

func TestSurfaceGradientAt(t *testing.T) {
	bound := NewBound(2, 0, 2, 0)
	surface := NewSurface(bound, 3, 3)

	surface.Grid[0] = []float64{0, 1, 2}
	surface.Grid[1] = []float64{5, 4, 3}
	surface.Grid[2] = []float64{7, 7, 7}

	// super simple octave code
	/*
		grid = [0, 1, 2; 5, 4, 3; 7, 7, 7]'
		delta = 0.01;

		points = [0.0, 0.0; 1.0, 1.0; 2.0, 2.0;...
		0.5, 0.0; 0.5, 0.75; 1.5, 1.75; 1.25, 1.25;...
		0.75, 1.25; 0.3, 1.0; 0.0, 1.25; 1.25, 0.0;...
		1.25, 2.0; 2.0, .75];

		for i = 1:max(size(points))
			d = max(size(grid))-1;
			xy = 0:d;

			if (points(i, 1) == d)
				points(i, 1) = points(i, 1) - delta;
			endif

			if (points(i, 2) == d)
				points(i, 2) = points(i, 2) - delta;
			endif

			if (points(i, 1) == 0)
				points(i, 1) = points(i, 1) + delta;
			endif

			if (points(i, 2) == 0)
				points(i, 2) = points(i, 2) + delta;
			endif

			x1 = interp2(xy, xy, grid, points(i, 1)-delta, points(i, 2));
			x2 = interp2(xy, xy, grid, points(i, 1)+delta, points(i, 2));

			y1 = interp2(xy, xy, grid, points(i, 1), points(i, 2)-delta);
			y2 = interp2(xy, xy, grid, points(i, 1), points(i, 2)+delta);

			[(x2-x1) / (2*delta), (y2 - y1) / (2*delta)]
		end
	*/

	tests := []testPointPoint{
		{&Point{0, 0}, &Point{5.0, 1.0}},
		{&Point{1, 1}, &Point{3.0, -1.0}},
		{&Point{2, 2}, &Point{4.0, 0.0}},
		{&Point{0.5, 0}, &Point{5.0, 0.0}},
		{&Point{0.5, 0.75}, &Point{3.5, 0.0}},
		{&Point{1.5, 1.75}, &Point{3.75, -0.5}},
		{&Point{1.25, 1.25}, &Point{3.25, -0.75}},
		{&Point{0.75, 1.25}, &Point{2.5, -0.5}},
		{&Point{0.3, 1.0}, &Point{3.0, 0.4}},
		{&Point{0.0, 1.25}, &Point{2.5, 1.0}},
		{&Point{1.25, 0.0}, &Point{2.0, -0.75}},
		{&Point{1.25, 2.0}, &Point{4.0, -0.75}},
		{&Point{2.0, 0.75}, &Point{2.75, 0.0}},
		{&Point{10.0, 0.75}, &Point{0.0, 0.0}},
	}

	for i, point := range tests {
		if v := surface.GradientAt(point.P); areaPointsDifferent(v, point.A, epsilon) {
			t.Errorf("surface, (%d) incorrect gradient at: expected %v, got %v", i, point.A, v)
		}
	}
}

func TestSurfaceWriteOffFile(t *testing.T) {
	bound := NewBound(3, 0, 3, 0)
	surface := NewSurface(bound, 4, 4)

	surface.Grid[1][1] = 0
	surface.Grid[2][1] = 1
	surface.Grid[1][2] = 2
	surface.Grid[2][2] = 3

	expected := "OFF\n16 5 0\n0.00000000 0.00000000 0.00000000\n0.00000000 1.00000000 0.00000000\n0.00000000 2.00000000 0.00000000\n0.00000000 3.00000000 0.00000000\n"
	expected += "1.00000000 0.00000000 0.00000000\n1.00000000 1.00000000 1.00000000\n1.00000000 2.00000000 3.00000000\n1.00000000 3.00000000 0.00000000\n"
	expected += "2.00000000 0.00000000 0.00000000\n2.00000000 1.00000000 0.00000000\n2.00000000 2.00000000 2.00000000\n2.00000000 3.00000000 0.00000000\n"
	expected += "3.00000000 0.00000000 0.00000000\n3.00000000 1.00000000 0.00000000\n3.00000000 2.00000000 0.00000000\n3.00000000 3.00000000 0.00000000\n"
	expected += "4 0 1 5 4\n4 2 3 7 6\n4 5 6 10 9\n4 8 9 13 12\n4 10 11 15 14\n"

	result := bytes.NewBufferString("")
	surface.WriteOffFile(result)

	if result.String() != expected {
		t.Errorf("surface, writeOffFile not right, %v != %v", result.String(), expected)
	}
}

func areaPointsDifferent(a, b *Point, delta float64) bool {
	return math.Abs(a[0]-b[0]) > delta || math.Abs(a[1]-b[1]) > delta
}
