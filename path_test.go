package geo

import (
	"bytes"
	"math"
	"math/rand"
	"testing"
)

func TestNewPathPreallocate(t *testing.T) {
	p := NewPathPreallocate(10, 1000)
	if l := p.Length(); l != 10 {
		t.Errorf("path, length not set correctly, got %d", l)
	}

	if c := cap(p.Points()); c != 1000 {
		t.Errorf("path, capactity not set corrctly, got %d", c)
	}

	p = NewPathPreallocate(100, 10)
	if l := p.Length(); l != 100 {
		t.Error("path, should handle length > capacity")
	}
}

func TestNewPathFromEncoding(t *testing.T) {
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

func TestNewPathFromEncodingShouldFailedOnMalformedPolyline(t *testing.T) {
	path := NewPathFromEncoding("xxxxxx")
	if path.Length() > 0 {
		t.Errorf("path should be empty, %v", path)
	}
}

func TestNewPathFromXYData(t *testing.T) {
	data := [][2]float64{
		{1, 2},
		{3, 4},
	}

	p := NewPathFromXYData(data)
	if l := p.Length(); l != len(data) {
		t.Errorf("path, should take full length of data, expected %d, got %d", len(data), l)
	}

	if point := p.GetAt(0); !point.Equals(&Point{1, 2}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}

	if point := p.GetAt(1); !point.Equals(&Point{3, 4}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}
}

func TestNewPathFromYXData(t *testing.T) {
	data := [][2]float64{
		{1, 2},
		{3, 4},
	}

	p := NewPathFromYXData(data)
	if l := p.Length(); l != len(data) {
		t.Errorf("path, should take full length of data, expected %d, got %d", len(data), l)
	}

	if point := p.GetAt(0); !point.Equals(&Point{2, 1}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}

	if point := p.GetAt(1); !point.Equals(&Point{4, 3}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}
}

func TestNewPathFromXYSlice(t *testing.T) {
	data := [][]float64{
		{1, 2, -1},
		nil,
		{3, 4},
	}

	p := NewPathFromXYSlice(data)
	if l := p.Length(); l != 2 {
		t.Errorf("path, should take full length of data, expected %d, got %d", 2, l)
	}

	if point := p.GetAt(0); !point.Equals(&Point{1, 2}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}

	if point := p.GetAt(1); !point.Equals(&Point{3, 4}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}
}

func TestNewPathFromYXSlice(t *testing.T) {
	data := [][]float64{
		{1, 2},
		{3, 4, -1},
	}

	p := NewPathFromYXSlice(data)
	if l := p.Length(); l != len(data) {
		t.Errorf("path, should take full length of data, expected %d, got %d", len(data), l)
	}

	if point := p.GetAt(0); !point.Equals(&Point{2, 1}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}

	if point := p.GetAt(1); !point.Equals(&Point{4, 3}) {
		t.Errorf("path, first point incorrect, got %v", point)
	}
}

func TestPathSetPoints(t *testing.T) {
	p := NewPath()

	points := make([]Point, 3)
	points[0] = *NewPoint(0, 0)
	points[1] = *NewPoint(1, 1)
	points[1] = *NewPoint(2, 2)

	p.SetPoints(points)
	if p.Length() != 3 {
		t.Error("path, set point length doesn't match")
	}
}

func TestPathPoints(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(0.5, .2))
	p.Push(NewPoint(1, 0))

	points := p.Points()
	if len(points) != 3 {
		t.Error("path, get point length doesn't match")
	}

	expected := NewPoint(0.5, 0.2)
	if !points[1].Equals(expected) {
		t.Errorf("path, get point points not equal, expected %v, got %v", expected, points[1])
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

	// empty path
	path := NewPath()
	if path.Encode() != "" {
		t.Error("path, encode empty path should be empty string")
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

func TestPathSquaredDistanceFrom(t *testing.T) {
	var answer float64

	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(0, 3))
	p.Push(NewPoint(4, 3))
	p.Push(NewPoint(4, 0))

	answer = 0.25
	if d := p.SquaredDistanceFrom(NewPoint(4.5, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, squaredDistanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.16
	if d := p.SquaredDistanceFrom(NewPoint(0.4, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, squaredDistanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.09
	if d := p.SquaredDistanceFrom(NewPoint(-0.3, 1.5)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, squaredDistanceFrom expected %f, got: %f", answer, d)
	}

	answer = 0.04
	if d := p.SquaredDistanceFrom(NewPoint(0.3, 2.8)); math.Abs(d-answer) > epsilon {
		t.Errorf("path, squaredDistanceFrom expected %f, got: %f", answer, d)
	}
}

func TestDirectionAt(t *testing.T) {
	path := NewPath().
		Push(NewPoint(0, 0)).
		Push(NewPoint(0, 1)).
		Push(NewPoint(1, 1)).
		Push(NewPoint(1, 0))

	// uses two surrounding points so directions are diagonal
	answers := []float64{0.5 * math.Pi, 0.25 * math.Pi, -0.25 * math.Pi}
	for i, v := range answers {
		if d := path.DirectionAt(i); d != v {
			t.Errorf("path, directionAt, expected %f, got %f", v, d)
		}
	}

	// INF for single point paths
	path = NewPath().Push(NewPoint(0, 0))
	if d := path.DirectionAt(0); d != math.Inf(1) {
		t.Errorf("path, directionAt expected Inf, got %f", d)
	}
}

func TestPathDirectionAtPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("path, directionAt invalid geometry should panic")
		}
	}()

	// these should panic
	NewPath().DirectionAt(0)
}

func TestPathMeasure(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(6, 8))
	p.Push(NewPoint(12, 0))

	result := p.Measure(NewPoint(3, 4))
	expected := 5.0
	if result != expected {
		t.Errorf("path, measure expected %f, got %f", expected, result)
	}

	// coincident with start point
	result = p.Measure(NewPoint(0, 0))
	expected = 0.0
	if result != expected {
		t.Errorf("path, measure expected %f, got %f", expected, result)
	}

	// coincident with end point
	result = p.Measure(NewPoint(12, 0))
	expected = 20.0
	if result != expected {
		t.Errorf("path, measure expected %f, got %f", expected, result)
	}

	// closest point on path
	result = p.Measure(NewPoint(-1, -1))
	expected = 0.0
	if result != expected {
		t.Errorf("path, measure expected %f, got %f", expected, result)
	}
}

func TestPathInterpolate(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(1, 1))
	p.Push(NewPoint(2, 2))
	p.Push(NewPoint(3, 3))

	// out-of-range - percent too low
	result := p.Interpolate(-0.1)
	expected := p.First()
	if !expected.Equals(result) {
		t.Errorf("path, interpolate expected %v, got %v", expected, result)
	}

	// start
	result = p.Interpolate(0)
	expected = p.First()
	if !expected.Equals(result) {
		t.Errorf("path, interpolate expected %v, got %v", expected, result)
	}

	// quarter
	result = p.Interpolate(0.25)
	expected = NewPoint(0.75, 0.75)
	if !expected.Equals(result) {
		t.Errorf("path, interpolate expected %v, got %v", expected, result)
	}

	// half
	result = p.Interpolate(0.5)
	expected = NewPoint(1.50, 1.50)
	if !expected.Equals(result) {
		t.Errorf("path, interpolate expected %v, got %v", expected, result)
	}

	// three quarters
	result = p.Interpolate(0.75)
	expected = NewPoint(2.25, 2.25)
	if !expected.Equals(result) {
		t.Errorf("path, interpolate expected %v, got %v", expected, result)
	}

	// end
	result = p.Interpolate(1)
	expected = p.Last()
	if !expected.Equals(result) {
		t.Errorf("path, interpolate expected %v, got %v", expected, result)
	}

	// out-of-range - percent too high
	result = p.Interpolate(1.1)
	expected = p.Last()
	if !expected.Equals(result) {
		t.Errorf("path, interpolate expected %v, got %v", expected, result)
	}
}

func TestPathProject(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(6, 8))
	p.Push(NewPoint(12, 0))

	result := p.Project(NewPoint(3, 4))
	expected := 0.25
	if result != expected {
		t.Errorf("path, project expected %f, got %f", expected, result)
	}

	// closest to the start
	result = p.Project(NewPoint(-1, -1))
	expected = 0.0
	if result != expected {
		t.Errorf("path, project expected %f, got %f", expected, result)
	}

	// closest to the end
	result = p.Project(NewPoint(13, -1))
	expected = 1.0
	if result != expected {
		t.Errorf("path, project expected %f, got %f", expected, result)
	}
}

func TestPathIntersection(t *testing.T) {
	path := NewPath()

	// these shouldn't panic
	path.Intersection(NewPath())
	path.Intersection(*NewPath())

	path.Intersection(NewLine(NewPoint(0, 0), NewPoint(1, 1)))
	path.Intersection(*NewLine(NewPoint(0, 0), NewPoint(1, 1)))
}

func TestPathIntersectionPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("path, intersection invalid geometry should panic")
		}
	}()

	// these should panic
	NewPath().Intersection(NewPoint(0, 0))
}

func TestPathIntersectionPath(t *testing.T) {
	var path *Path
	var answer *Point

	p := NewPath().Push(NewPoint(0, 0)).Push(NewPoint(1, 1)).Push(NewPoint(2, 2))

	answer = NewPoint(0.5, 0.5)
	path = NewPath()
	path.Push(NewPoint(0, 0.5)).Push(NewPoint(1, 0.5))
	if p, i := p.IntersectionPath(path); !p[0].Equals(answer) || i[0][0] != 0 || i[0][1] != 0 || len(p) != 1 || len(i) != 1 {
		t.Errorf("path, intersectionPath expected %v, got: %v, %v", answer, p, i)
	}

	answer = NewPoint(1.5, 1.5)
	path = NewPath()
	path.Push(NewPoint(0, 1.5)).Push(NewPoint(2, 1.5))
	if p, i := p.IntersectionPath(path); !p[0].Equals(answer) || i[0][0] != 1 || i[0][1] != 0 || len(p) != 1 || len(i) != 1 {
		t.Errorf("path, intersectionPath expected %v, got: %v, %v", answer, p, i)
	}

	answer = NewPoint(1.5, 1.5)
	path = NewPath()
	path.Push(NewPoint(0, 1.5)).Push(NewPoint(1, 1.5)).Push(NewPoint(2, 1.5))
	if p, i := p.IntersectionPath(path); !p[0].Equals(answer) || i[0][0] != 1 || i[0][1] != 1 || len(p) != 1 || len(i) != 1 {
		t.Errorf("path, intersectionPath expected %v, got: %v, %v", answer, p, i)
	}

	path = NewPath()
	path.Push(NewPoint(0, 1.5)).Push(NewPoint(1, 1.5))
	if p, i := p.IntersectionPath(path); len(p) != 0 || len(i) != 0 {
		t.Errorf("path, intersectionPath expected none, got: %v, %v", p, i)
	}
}

func TestPathIntersectionLine(t *testing.T) {
	var line *Line
	var answer *Point

	p := NewPath().Push(NewPoint(0, 0)).Push(NewPoint(1, 1)).Push(NewPoint(2, 2))

	answer = NewPoint(0.5, 0.5)
	line = NewLine(NewPoint(0, 0.5), NewPoint(1, 0.5))
	if p, i := p.IntersectionLine(line); !p[0].Equals(answer) || i[0][0] != 0 || i[0][1] != 0 || len(p) != 1 || len(i) != 1 {
		t.Errorf("path, intersectionLine expected %v, got: %v, %v", answer, p, i)
	}

	answer = NewPoint(1.5, 1.5)
	line = NewLine(NewPoint(0, 1.5), NewPoint(2, 1.5))
	if p, i := p.IntersectionLine(line); !p[0].Equals(answer) || i[0][0] != 1 || i[0][1] != 0 || len(p) != 1 || len(i) != 1 {
		t.Errorf("path, intersectionLine expected %v, got: %v, %v", answer, p, i)
	}

	line = NewLine(NewPoint(0, 1.5), NewPoint(1, 1.5))
	if p, i := p.IntersectionLine(line); len(p) != 0 || len(i) != 0 {
		t.Errorf("path, intersectionLine expected none, got: %v, %v", p, i)
	}
}

func TestPathIntersects(t *testing.T) {
	path := NewPath()

	// these shouldn't panic
	path.Intersects(NewPath())
	path.Intersects(*NewPath())

	path.Intersects(NewLine(NewPoint(0, 0), NewPoint(1, 1)))
	path.Intersects(*NewLine(NewPoint(0, 0), NewPoint(1, 1)))
}

func TestPathIntersectsPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("path, intersects invalid geometry should panic")
		}
	}()

	// these should panic
	NewPath().Intersects(NewPoint(0, 0))
}

func TestPathIntersectsPath(t *testing.T) {
	var path *Path
	var answer bool

	p := NewPath().Push(NewPoint(0, 0)).Push(NewPoint(1, 1)).Push(NewPoint(2, 2))

	answer = true
	path = NewPath()
	path.Push(NewPoint(0, 0.5)).Push(NewPoint(1, 0.5))
	if b := p.IntersectsPath(path); b != answer {
		t.Errorf("path, intersectsPath expected %v, got: %v", answer, b)
	}

	answer = true
	path = NewPath()
	path.Push(NewPoint(0, 1)).Push(NewPoint(1, 1))
	if b := p.IntersectsPath(path); b != answer {
		t.Errorf("path, intersectsPath expected %v, got: %v", answer, b)
	}

	answer = false
	path = NewPath()
	path.Push(NewPoint(0, 1)).Push(NewPoint(0, 2))
	if b := p.IntersectsPath(path); b != answer {
		t.Errorf("path, intersectsPath expected %v, got: %v", answer, b)
	}
}

func TestPathIntersectsLine(t *testing.T) {
	var line *Line
	var answer bool

	p := NewPath().Push(NewPoint(0, 0)).Push(NewPoint(1, 1)).Push(NewPoint(2, 2))

	answer = true
	line = NewLine(NewPoint(0, 0.5), NewPoint(1, 0.5))
	if b := p.IntersectsLine(line); b != answer {
		t.Errorf("path, intersectsLine expected %v, got: %v", answer, b)
	}

	answer = true
	line = NewLine(NewPoint(0, 1), NewPoint(1, 1))
	if b := p.IntersectsLine(line); b != answer {
		t.Errorf("path, intersectsLine expected %v, got: %v", answer, b)
	}

	answer = false
	line = NewLine(NewPoint(0, 1), NewPoint(0, 2))
	if b := p.IntersectsLine(line); b != answer {
		t.Errorf("path, intersectsLine expected %v, got: %v", answer, b)
	}
}

func TestPathWriteOffFile(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(0.5, .2))
	p.Push(NewPoint(1, 0))

	expected := "OFF\n3 1 0\n0.000000 0.000000 0\n0.500000 0.200000 0\n1.000000 0.000000 0\n3 0 1 2 170 170 170\n"
	result := bytes.NewBufferString("")
	p.WriteOffFile(result)

	if off := result.String(); off != expected {
		t.Errorf("path, writeOffFile not right, %v != %v", expected, off)
	}

	expected = "OFF\n3 1 0\n0.000000 0.000000 0\n0.500000 0.200000 0\n1.000000 0.000000 0\n3 0 1 2 1 2 3\n"
	result = bytes.NewBufferString("")
	p.WriteOffFile(result, [3]int{1, 2, 3})

	if off := result.String(); off != expected {
		t.Errorf("path, writeOffFile not right, %v != %v", expected, off)
	}
}

func TestPathToGeoJSON(t *testing.T) {
	p := NewPath().
		Push(NewPoint(1, 2))

	f := p.ToGeoJSON()
	if !f.Geometry.IsLineString() {
		t.Errorf("path, should be linestring geometry")
	}
}

func TestPathToWKT(t *testing.T) {
	p := NewPath()

	answer := "EMPTY"
	if s := p.ToWKT(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}

	p.Push(NewPoint(1, 2))
	answer = "LINESTRING(1 2)"
	if s := p.ToWKT(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}

	p.Push(NewPoint(3, 4))
	answer = "LINESTRING(1 2,3 4)"
	if s := p.ToWKT(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}
}

func TestPathString(t *testing.T) {
	p := NewPath()

	answer := "EMPTY"
	if s := p.String(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}

	p.Push(NewPoint(1, 2))
	answer = "LINESTRING(1 2)"
	if s := p.String(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}

	p.Push(NewPoint(3, 4))
	answer = "LINESTRING(1 2,3 4)"
	if s := p.String(); s != answer {
		t.Errorf("path, string expected %s, got %s", answer, s)
	}
}
