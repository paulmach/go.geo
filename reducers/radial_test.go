package reducers

import (
	"reflect"
	"testing"

	"github.com/paulmach/go.geo"
)

func TestRadial(t *testing.T) {
	p := geo.NewPath()
	if len(Radial(p, 1.0)) != 0 {
		t.Error("radial could not reduce zero length path")
	}

	p = append(p,
		geo.NewPoint(0, 0),
		geo.NewPoint(0, 1),
		geo.NewPoint(0, 2),
	)

	if l := len(Radial(p, 0.9)); l != 3 {
		t.Error("radial reduce to incorrect number of points")
	}

	reduced := Radial(p, 1.1)
	if l := len(reduced); l != 2 {
		t.Error("radial reduce to incorrect number of points")
	}

	p = append(p,
		geo.NewPoint(0, 3),
		geo.NewPoint(0, 4),
		geo.NewPoint(0, 5),
	)

	if l := len(Radial(p, 1.1)); l != 4 {
		t.Error("radial reduce to incorrect number of points")
	}
}

func TestRadialIndexMap(t *testing.T) {
	p := geo.NewPath()
	if reduced, _ := RadialIndexMap(p, 1.0); len(reduced) != 0 {
		t.Error("radialIndexMap could not reduce zero length path")
	}

	p = append(p,
		geo.NewPoint(0, 0),
		geo.NewPoint(0, 1),
		geo.NewPoint(0, 2),
	)

	reduced, im := RadialIndexMap(p, 0.9)
	if len(reduced) != 3 {
		t.Error("radialIndexMap reduce to incorrect number of points")
	}

	if !reflect.DeepEqual(im, []int{0, 1, 2}) {
		t.Errorf("radialIndexMap reduce bad index map, got %v", im)
	}

	reduced, im = RadialIndexMap(p, 1.1)
	if l := len(reduced); l != 2 {
		t.Error("radialIndexMap reduce to incorrect number of points")
	}

	if !reflect.DeepEqual(im, []int{0, 2}) {
		t.Errorf("radialIndexMap reduce bad index map, got %v", im)
	}

	p = append(p,
		geo.NewPoint(0, 3),
		geo.NewPoint(0, 4),
		geo.NewPoint(0, 5),
	)

	reduced, im = RadialIndexMap(p, 1.1)
	if len(reduced) != 4 {
		t.Error("radialIndexMap reduce to incorrect number of points")
	}

	if !reflect.DeepEqual(im, []int{0, 2, 4, 5}) {
		t.Errorf("radialIndexMap reduce bad index map, got %v", im)
	}
}

func TestRadialGeo(t *testing.T) {
	p := append(geo.NewPath(),
		geo.NewPoint(0, 0),
		geo.NewPoint(0, 1),
		geo.NewPoint(0, 2),
	)

	threshold := 1.1
	if l := len(RadialGeo(p, threshold)); l != 3 {
		t.Error("radialGeo reduce to incorrect number of points")
	}

	threshold = p[0].GeoDistanceFrom(p[1]) + 1.0
	reduced := RadialGeo(p, threshold)
	if l := len(reduced); l != 2 {
		t.Error("radialGeo reduce to incorrect number of points")
	}
}

func TestRadialGeoIndexMap(t *testing.T) {
	p := append(geo.NewPath(),
		geo.NewPoint(0, 0),
		geo.NewPoint(0, 1),
		geo.NewPoint(0, 2),
	)

	threshold := 1.1
	reduced, im := RadialGeoIndexMap(p, threshold)
	if len(reduced) != 3 {
		t.Error("radialGeoIndexMap reduce to incorrect number of points")
	}

	if !reflect.DeepEqual(im, []int{0, 1, 2}) {
		t.Errorf("radialGeoIndexMap reduce bad index map, got %v", im)
	}

	threshold = p[0].GeoDistanceFrom(p[1]) + 1.0
	reduced, im = RadialGeoIndexMap(p, threshold)
	if l := len(reduced); l != 2 {
		t.Error("radialGeoIndexMap reduce to incorrect number of points")
	}

	if !reflect.DeepEqual(im, []int{0, 2}) {
		t.Errorf("radialGeoIndexMap reduce bad index map, got %v", im)
	}
}
