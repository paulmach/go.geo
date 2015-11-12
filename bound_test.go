package geo

import (
	"math"
	"testing"
)

func TestBoundNew(t *testing.T) {
	bound := NewBound(5, 0, 3, 0)
	if !bound.sw.Equals(NewPoint(0, 0)) {
		t.Errorf("bound, incorrect sw: expected %v, got %v", NewPoint(0, 0), bound.sw)
	}

	if !bound.ne.Equals(NewPoint(5, 3)) {
		t.Errorf("bound, incorrect ne: expected %v, got %v", NewPoint(5, 3), bound.ne)
	}

	bound = NewBoundFromPoints(NewPoint(0, 3), NewPoint(4, 0))
	if !bound.sw.Equals(NewPoint(0, 0)) {
		t.Errorf("bound, incorrect sw: expected %v, got %v", NewPoint(0, 0), bound.sw)
	}

	if !bound.ne.Equals(NewPoint(4, 3)) {
		t.Errorf("bound, incorrect ne: expected %v, got %v", NewPoint(4, 3), bound.ne)
	}

	bound1 := NewBound(1, 2, 3, 4)
	bound2 := NewBoundFromPoints(NewPoint(1, 3), NewPoint(2, 4))
	if !bound1.Equals(bound2) {
		t.Errorf("bound, expected %v == %v", bound1, bound2)
	}
}

func TestGeoBoundAroundPoint(t *testing.T) {
	p := &Point{}
	p.SetLat(50.0359)
	p.SetLng(5.42553)
	bound := NewGeoBoundAroundPoint(p, 1000000)
	if bound.Center().Lat() != p.Lat() {
		t.Errorf("bound, should have correct center lat point")
	}

	if bound.Center().Lng() != p.Lng() {
		t.Errorf("bound, should have correct center lng point")
	}

	//Given point is 968.9 km away from center
	if !bound.Contains(&Point{3.412, 58.3838}) {
		t.Errorf("bound, should have point included in bound")
	}

	bound = NewGeoBoundAroundPoint(p, 10000.0)
	if bound.Center().Lat() != p.Lat() {
		t.Errorf("bound, should have correct center lat point")
	}

	if bound.Center().Lng() != p.Lng() {
		t.Errorf("bound, should have correct center lng point")
	}

	//Given point is 968.9 km away from center
	if bound.Contains(&Point{3.412, 58.3838}) {
		t.Errorf("bound, should not have point included in bound")
	}
}

func TestNewBoundFromMapTile(t *testing.T) {
	bound := NewBoundFromMapTile(7, 8, 9)

	level := uint64(9 + 5) // we're testing point +5 zoom, in same tile
	factor := uint64(5)

	// edges should be within the bounds
	lng, lat := ScalarMercator.Inverse(7<<factor+1, 8<<factor+1, level)
	if !bound.Contains(NewPoint(lng, lat)) {
		t.Errorf("bound, should contain point")
	}

	lng, lat = ScalarMercator.Inverse(7<<factor-1, 8<<factor-1, level)
	if bound.Contains(NewPoint(lng, lat)) {
		t.Errorf("bound, should not contain point")
	}

	lng, lat = ScalarMercator.Inverse(8<<factor-1, 9<<factor-1, level)
	if !bound.Contains(NewPoint(lng, lat)) {
		t.Errorf("bound, should contain point")
	}

	lng, lat = ScalarMercator.Inverse(8<<factor+1, 9<<factor+1, level)
	if bound.Contains(NewPoint(lng, lat)) {
		t.Errorf("bound, should not contain point")
	}

	bound = NewBoundFromMapTile(7, 8, 35)
}

func TestBoundExtend(t *testing.T) {
	bound := NewBound(3, 0, 5, 0)

	if b := bound.Clone().Extend(NewPoint(2, 1)); !b.Equals(bound) {
		t.Errorf("bound, extend expected %v, got %v", bound, b)
	}

	answer := NewBound(6, 0, 5, -1)
	if b := bound.Clone().Extend(NewPoint(6, -1)); !b.Equals(answer) {
		t.Errorf("bound, extend expected %v, got %v", answer, b)
	}
}

func TestBoundUnion(t *testing.T) {
	b1 := NewBound(0, 1, 0, 1)
	b2 := NewBound(0, 2, 0, 2)

	expected := NewBound(0, 2, 0, 2)
	if b := b1.Clone().Union(b2); !b.Equals(expected) {
		t.Errorf("bound, expected %v, got %v", expected, b)
	}

	if b := b2.Clone().Union(b1); !b.Equals(expected) {
		t.Errorf("bound, expected %v, got %v", expected, b)
	}
}

func TestBoundContains(t *testing.T) {
	var p *Point
	bound := NewBound(2, -2, 1, -1)

	p = NewPoint(0, 0)
	if !bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to be within %v", p, bound)
	}

	p = NewPoint(-1, 0)
	if !bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to be within %v", p, bound)
	}

	p = NewPoint(2, 1)
	if !bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to be within %v", p, bound)
	}

	p = NewPoint(0, 3)
	if bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to not be within %v", p, bound)
	}

	p = NewPoint(0, -3)
	if bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to not be within %v", p, bound)
	}

	p = NewPoint(3, 0)
	if bound.Contains(p) {
		t.Errorf("bound, contains expected %v, to not be within %v", p, bound)
	}

	p = NewPoint(-3, 0)
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

	a := NewBound(7, 8, 6, 7)
	b := NewBound(6.1, 8.1, 6.1, 8.1)

	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("bound, expected to intersect")
	}

	a = NewBound(1, 4, 2, 3)
	b = NewBound(2, 3, 1, 4)

	if !a.Intersects(b) || !b.Intersects(a) {
		t.Errorf("bound, expected to intersect")
	}
}

func TestBoundCenter(t *testing.T) {
	var p *Point
	var b *Bound

	b = NewBound(0, 1, 2, 3)
	p = NewPoint(0.5, 2.5)
	if c := b.Center(); !c.Equals(p) {
		t.Errorf("bound, center expected %v, got %v", p, c)
	}

	b = NewBound(0, 0, 2, 2)
	p = NewPoint(0, 2)
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

func TestBoundGeoPad(t *testing.T) {
	tests := []*Bound{
		NewBoundFromPoints(NewPoint(-122.559, 37.887), NewPoint(-122.521, 37.911)),
		NewBoundFromPoints(NewPoint(-122.559, 15), NewPoint(-122.521, 15)),
		NewBoundFromPoints(NewPoint(20, -15), NewPoint(20, -15)),
	}

	for i, b1 := range tests {
		b2 := b1.Clone().GeoPad(100)

		if math.Abs(b1.GeoHeight()+200-b2.GeoHeight()) > 1.0 {
			t.Errorf("bound, geoPad height incorrected for %d, expected %v, got %v", i, b1.GeoHeight()+200, b2.GeoHeight())
		}

		if math.Abs(b1.GeoWidth()+200-b2.GeoWidth()) > 1.0 {
			t.Errorf("bound, geoPad width incorrected for %d, expected %v, got %v", i, b1.GeoWidth()+200, b2.GeoWidth())
		}
	}
}

func TestBoundAccessors(t *testing.T) {
	bound := NewBound(1, 2, 3, 4)

	if !bound.sw.Equals(bound.SouthWest()) || !bound.SouthWest().Equals(bound.sw) {
		t.Errorf("bound, southwest expected %v == %v", bound.sw, bound.SouthWest())
	}

	if !bound.ne.Equals(bound.NorthEast()) || !bound.NorthEast().Equals(bound.ne) {
		t.Errorf("bound, northeast expected %v == %v", bound.ne, bound.NorthEast())
	}

	if !bound.NorthWest().Equals(NewPoint(1, 4)) {
		t.Errorf("bound, northwest incorrect, got %v", bound.NorthWest())
	}

	if !bound.SouthEast().Equals(NewPoint(2, 3)) {
		t.Errorf("bound, southeast incorrect, got %v", bound.SouthEast())
	}
}

func TestBoundEquals(t *testing.T) {
	bound1 := NewBound(1, 2, 3, 4)
	bound2 := NewBoundFromPoints(NewPoint(1, 3), NewPoint(2, 4))
	if !bound1.Equals(bound2) || !bound2.Equals(bound1) {
		t.Errorf("bound, expected %v == %v", bound1, bound2)
	}

	bound2 = NewBound(1, 2, 4, 4)
	if bound1.Equals(bound2) || bound2.Equals(bound1) {
		t.Errorf("bound, expected %v != %v", bound1, bound2)
	}

	bound2 = NewBound(1, 1, 3, 4)
	if bound1.Equals(bound2) || bound2.Equals(bound1) {
		t.Errorf("bound, expected %v != %v", bound1, bound2)
	}
}

func TestBoundSides(t *testing.T) {
	// NewBound(west, east, south, north)
	b := NewBound(1, 2, 3, 4)

	if v := b.North(); v != 4 {
		t.Errorf("incorrect north value, got %v", v)
	}

	if v := b.South(); v != 3 {
		t.Errorf("incorrect south value, got %v", v)
	}

	if v := b.East(); v != 2 {
		t.Errorf("incorrect east value, got %v", v)
	}

	if v := b.West(); v != 1 {
		t.Errorf("incorrect west value, got %v", v)
	}

	if b.SouthWest().Lng() != b.West() {
		t.Errorf("west value incorrect")
	}

	if b.SouthWest().Lat() != b.South() {
		t.Errorf("south value incorrect")
	}

	if b.NorthEast().Lng() != b.East() {
		t.Errorf("east value incorrect")
	}

	if b.NorthEast().Lat() != b.North() {
		t.Errorf("north value incorrect")
	}

	if b.SouthWest().Lng() != b.Left() {
		t.Errorf("left value incorrect")
	}

	if b.SouthWest().Lat() != b.Bottom() {
		t.Errorf("bottom value incorrect")
	}

	if b.NorthEast().Lng() != b.Right() {
		t.Errorf("right value incorrect")
	}

	if b.NorthEast().Lat() != b.Top() {
		t.Errorf("top value incorrect")
	}
}

func TestBoundEmpty(t *testing.T) {
	bound := NewBound(1, 2, 3, 4)
	if bound.Empty() {
		t.Error("bound, empty exported false, got true")
	}

	bound = NewBound(1, 1, 2, 2)
	if !bound.Empty() {
		t.Error("bound, empty exported true, got false")
	}

	// horizontal bar
	bound = NewBound(1, 1, 2, 3)
	if !bound.Empty() {
		t.Error("bound, empty exported true, got false")
	}

	// vertical bar
	bound = NewBound(1, 2, 2, 2)
	if !bound.Empty() {
		t.Error("bound, empty exported true, got false")
	}

	// negative/malformed area
	bound = NewBound(1, 0, 2, 2)
	if !bound.Empty() {
		t.Error("bound, empty exported true, got false")
	}

	// negative/malformed area
	bound = NewBound(1, 1, 2, 1)
	if !bound.Empty() {
		t.Error("bound, empty exported true, got false")
	}
}

func TestBoundString(t *testing.T) {
	bound := NewBound(1, 2, 3, 4)

	answer := "POLYGON((1 3, 1 4, 2 4, 2 3, 1 3))"
	if s := bound.String(); s != answer {
		t.Errorf("bound, string expected %s, got %s", answer, s)
	}
}

func TestBoundToMysqlIntersectsCondition(t *testing.T) {
	b := NewBound(1, 2, 3, 4)

	expected := "INTERSECTS(column, GEOMFROMTEXT('POLYGON((1 3, 1 4, 2 4, 2 3, 1 3))'))"
	if p := b.ToMysqlIntersectsCondition("column"); p != expected {
		t.Errorf("bound, incorrect condition, got %v", p)
	}
}
