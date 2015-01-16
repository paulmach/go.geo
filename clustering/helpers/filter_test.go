package helpers

import (
	"testing"

	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/clustering"
)

func TestFilterSmallClusters(t *testing.T) {
	g := FilterSmallClusters([]*clustering.Cluster{}, 5)
	if g == nil {
		t.Errorf("result should not be nil")
	}

	if len(g) != 0 {
		t.Errorf("zero group input should result in zero group output")
	}

	testSet := []*clustering.Cluster{
		clustering.NewCluster(&event{Location: geo.NewPoint(1, 1)}),
		clustering.NewCluster(&event{Location: geo.NewPoint(1, 1)}, &event{Location: geo.NewPoint(2, 2)}),
	}

	g = FilterSmallClusters(testSet, 5)

	if l := len(g); l != 0 {
		t.Errorf("should filter out small groups, but got %d", l)
	}

	if l := len(testSet); l != 2 {
		t.Errorf("should not change test set, but got length %d", l)
	}

	g = FilterSmallClusters(testSet, 2)
	if l := len(g); l != 1 {
		t.Errorf("should filter out small groups, but got %d", l)
	}

	if l := len(testSet); l != 2 {
		t.Errorf("should not change test set, but got length %d", l)
	}
}

type event struct {
	Location *geo.Point
}

func (e *event) CenterPoint() *geo.Point {
	return e.Location
}
