package reducers

import (
	"reflect"
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

func TestRadialIndexMap(t *testing.T) {
	p := geo.NewPath()
	if reduced, _ := RadialIndexMap(p, 1.0); reduced.Length() != 0 {
		t.Error("radialIndexMap could not reduce zero length path")
	}

	p.Push(geo.NewPoint(0, 0))
	p.Push(geo.NewPoint(0, 1))
	p.Push(geo.NewPoint(0, 2))

	reduced, im := RadialIndexMap(p, 0.9)
	if reduced.Length() != 3 {
		t.Error("radialIndexMap reduce to incorrect number of points")
	}

	if !reflect.DeepEqual(im, []int{0, 1, 2}) {
		t.Errorf("radialIndexMap reduce bad index map, got %v", im)
	}

	reduced, im = RadialIndexMap(p, 1.1)
	if l := reduced.Length(); l != 2 {
		t.Error("radialIndexMap reduce to incorrect number of points")
	}

	if !reflect.DeepEqual(im, []int{0, 2}) {
		t.Errorf("radialIndexMap reduce bad index map, got %v", im)
	}

	if reduced == p {
		t.Error("should create new path and not modify original")
	}

	p.Push(geo.NewPoint(0, 3))
	p.Push(geo.NewPoint(0, 4))
	p.Push(geo.NewPoint(0, 5))

	reduced, im = RadialIndexMap(p, 1.1)
	if reduced.Length() != 4 {
		t.Error("radialIndexMap reduce to incorrect number of points")
	}

	if !reflect.DeepEqual(im, []int{0, 2, 4, 5}) {
		t.Errorf("radialIndexMap reduce bad index map, got %v", im)
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

func TestRadialGeoIndexMap(t *testing.T) {
	p := geo.NewPath()
	p.Push(geo.NewPoint(0, 0))
	p.Push(geo.NewPoint(0, 1))
	p.Push(geo.NewPoint(0, 2))

	threshold := 1.1
	reduced, im := RadialGeoIndexMap(p, threshold)
	if reduced.Length() != 3 {
		t.Error("radialGeoIndexMap reduce to incorrect number of points")
	}

	if !reflect.DeepEqual(im, []int{0, 1, 2}) {
		t.Errorf("radialGeoIndexMap reduce bad index map, got %v", im)
	}

	threshold = p.GetAt(0).GeoDistanceFrom(p.GetAt(1)) + 1.0
	reduced, im = RadialGeoIndexMap(p, threshold)
	if l := reduced.Length(); l != 2 {
		t.Error("radialGeoIndexMap reduce to incorrect number of points")
	}

	if !reflect.DeepEqual(im, []int{0, 2}) {
		t.Errorf("radialGeoIndexMap reduce bad index map, got %v", im)
	}

	if reduced == p {
		t.Error("should create new path and not modify original")
	}
}
