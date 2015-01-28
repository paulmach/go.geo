package geo

import (
	"math"
	"testing"
)

func TestGepCentroid(t *testing.T) {
	ps := &PointSet{}
	ps.Push(&Point{-122.42558918, 37.76159786}).
		Push(&Point{-122.41486043, 37.78138826}).
		Push(&Point{-122.40206146, 37.77962363})
	centroid := ps.GeoCentroid()
	expectedCenter := &Point{-122.41417035666666, 37.77420325}
	if !centroid.Equals(expectedCenter) {
		t.Errorf("should find centroid correctly, got %v", centroid.Lng())
	}
}

func TestGeoDistanceFrom(t *testing.T) {
	ps := &PointSet{}
	ps.Push(&Point{-122.42558918, 37.76159786}).
		Push(&Point{-122.41486043, 37.78138826}).
		Push(&Point{-122.40206146, 37.77962363})
	fromPoint := &Point{-122.41941550000001, 37.7749295}

	if distance := ps.GeoDistanceFrom(fromPoint); math.Floor(distance) != 823 {
		t.Errorf("should find geo distance from correctly, got %v", distance)
	}
}

func TestNewPointSet(t *testing.T) {
	ps := NewPointSet()
	ps.Push(&Point{-122.42558918, 37.76159786}).
		Push(&Point{-122.41486043, 37.78138826}).
		Push(&Point{-122.40206146, 37.77962363})
	if ps.Length() != 3 {
		t.Errorf("should find correct length of new point set %v", ps.Length())
	}
}

func TestNewPointSetPreallocate(t *testing.T) {
	ps := NewPointSet()
	ps.Push(&Point{-122.42558918, 37.76159786}).
		Push(&Point{-122.41486043, 37.78138826}).
		Push(&Point{-122.40206146, 37.77962363})

	if ps.Length() != 3 {
		t.Errorf("should find correct length of new point set %v", ps.Length())
	}

	if !ps.GetAt(0).Equals(&Point{-122.42558918, 37.76159786}) {
		t.Errorf("should find correct first point of new point set %v", ps.GetAt(0))
	}

	if !ps.GetAt(2).Equals(&Point{-122.40206146, 37.77962363}) {
		t.Errorf("should find correct first point of new point set %v", ps.GetAt(2))
	}

}
