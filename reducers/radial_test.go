package reducers

import (
	"testing"

	"github.com/paulmach/go.geo"
)

func TestRadial(t *testing.T) {
	p := geo.NewPath()
	if Radial(p, 1.0).Length() != 0 {
		t.Error("radial could not reduce zero length path")
	}

	p.Push(geo.NewPoint(0, 0))
	p.Push(geo.NewPoint(0, 1))
	p.Push(geo.NewPoint(0, 2))

	if l := Radial(p, 0.9).Length(); l != 3 {
		t.Error("radial reduce to incorrect number of points")
	}

	reduced := Radial(p, 1.1)
	if l := reduced.Length(); l != 2 {
		t.Error("radial reduce to incorrect number of points")
	}

	if reduced == p {
		t.Error("should create new path and not modify original")
	}

	p.Push(geo.NewPoint(0, 3))
	p.Push(geo.NewPoint(0, 4))
	p.Push(geo.NewPoint(0, 5))
	if l := Radial(p, 1.1).Length(); l != 4 {
		t.Error("radial reduce to incorrect number of points")
	}
}

func TestRadialGeo(t *testing.T) {
	p := geo.NewPath()
	p.Push(geo.NewPoint(0, 0))
	p.Push(geo.NewPoint(0, 1))
	p.Push(geo.NewPoint(0, 2))

	threshold := 1.1
	if l := RadialGeo(p, threshold).Length(); l != 3 {
		t.Error("radialGeo reduce to incorrect number of points")
	}

	threshold = p.GetAt(0).GeoDistanceFrom(p.GetAt(1)) + 1.0
	reduced := RadialGeo(p, threshold)
	if l := reduced.Length(); l != 2 {
		t.Error("radialGeo reduce to incorrect number of points")
	}

	if reduced == p {
		t.Error("should create new path and not modify original")
	}
}
