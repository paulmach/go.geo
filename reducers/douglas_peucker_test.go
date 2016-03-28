package reducers

import (
	"reflect"
	"testing"

	"github.com/paulmach/go.geo"
)

func TestDouglasPeucker(t *testing.T) {
	p := geo.NewPath()
	if reducedPath := DouglasPeucker(p, 0.1); !reducedPath.Equal(p) {
		t.Error("dp should return same path if of length 0")
	}

	p = append(p, geo.NewPoint(0, 0))
	if reducedPath := DouglasPeucker(p, 0.1); !reducedPath.Equal(p) {
		t.Error("dp should return same path if of length 1")
	}

	p = append(p, geo.NewPoint(0.5, .2))
	if reducedPath := DouglasPeucker(p, 0.1); !reducedPath.Equal(p) {
		t.Error("dp should return same path if of length 2")
	}

	p = append(p, geo.NewPoint(1, 0))

	if l := len(DouglasPeucker(p, 0.1)); l != 3 {
		t.Errorf("dp reduce to incorrect number of points, expected 2, got %d", l)
	}

	if l := len(DouglasPeucker(p, 0.3)); l != 2 {
		t.Errorf("dp reduce to incorrect number of points, expected 3, got %d", l)
	}

	p = append(geo.NewPath(),
		geo.NewPoint(0, 0),
		geo.NewPoint(0, 1),
		geo.NewPoint(0, 2),
	)

	if l := len(DouglasPeucker(p, 0.0)); l != 2 {
		t.Errorf("dp reduce should remove coliniar points")
	}
}

func TestDouglasPeuckerIndexMap(t *testing.T) {
	p := geo.NewPath()

	// zero length
	reduced, indexMap := DouglasPeuckerIndexMap(p, 0.1)
	if len(reduced) != 0 {
		t.Error("dpim should return same path if of length 0")
	}

	if len(indexMap) != 0 {
		t.Error("dpim should have map of zero length for empty path input")
	}

	// 1 length
	p = append(p, geo.NewPoint(0, 0))
	reduced, indexMap = DouglasPeuckerIndexMap(p, 0.1)
	if !reduced.Equal(p) {
		t.Error("dpim should return same path if of length 1")
	}

	if !reflect.DeepEqual(indexMap, []int{0}) {
		t.Error("dpim should return []int{0} for index map when path is length 1")
	}

	// 2 length
	p = append(p, geo.NewPoint(0.5, .2))
	reduced, indexMap = DouglasPeuckerIndexMap(p, 0.1)
	if !reduced.Equal(p) {
		t.Error("dpim should return same path if of length 2")
	}

	if !reflect.DeepEqual(indexMap, []int{0, 1}) {
		t.Error("dpim should return []int{0, 1} for index map when path is length 2")
	}

	// 3 length, does reduce
	p = append(p, geo.NewPoint(1, 0))
	reduced, indexMap = DouglasPeuckerIndexMap(p, 0.3)
	if l := len(reduced); l != 2 {
		t.Errorf("dpim reduce to incorrect number of points, expected 2, got %d", l)
	}

	if !reflect.DeepEqual(indexMap, []int{0, 2}) {
		t.Errorf("dpim should return []int{0, 2} for index map, got %v %v", indexMap, []int{0, 2})
	}

	// 3 length, doesn't reduce
	reduced, indexMap = DouglasPeuckerIndexMap(p, 0.1)
	if l := len(reduced); l != 3 {
		t.Errorf("dpim reduce to incorrect number of points, expected 3, got %d", l)
	}

	if !reflect.DeepEqual(indexMap, []int{0, 1, 2}) {
		t.Errorf("dpim should return []int{0, 1, 2} for index map, got %v", indexMap)
	}
}
