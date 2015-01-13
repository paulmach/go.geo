package helpers

import (
	"testing"

	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/clustering/point_clustering"
)

func TestRematchPointersToClusters(t *testing.T) {
	c := RematchPointersToClusters([]*point_clustering.Cluster{}, []point_clustering.Pointer{}, point_clustering.CentroidGeoDistance{}, 30)
	if c == nil {
		t.Errorf("result should not be nil")
	}

	if len(c) != 0 {
		t.Errorf("zero cluster input should result in zero cluster output")
	}

	testClusters := []*point_clustering.Cluster{
		point_clustering.NewClusterWithCentroid(geo.NewPoint(1, 1)),
		point_clustering.NewClusterWithCentroid(geo.NewPoint(2, 2)),
	}

	testPointers := []point_clustering.Pointer{
		&event{Location: geo.NewPoint(1, 1)},
		&event{Location: geo.NewPoint(1, 1)},
		&event{Location: geo.NewPoint(2, 2)},
		&event{Location: geo.NewPoint(3, 3)},
	}

	c = RematchPointersToClusters(testClusters, testPointers, point_clustering.CentroidDistance{}, 1)
	if l := len(c[0].Pointers); l != 2 {
		t.Errorf("wrong number of pointers, got %d", l)
	}

	if l := len(c[1].Pointers); l != 1 {
		t.Errorf("wrong number of pointers, got %d", l)
	}

	c = RematchPointersToClusters(testClusters, testPointers, point_clustering.CentroidDistance{}, 2)

	if l := len(c[0].Pointers); l != 2 {
		t.Errorf("wrong number of pointers, got %d", l)
	}

	if l := len(c[1].Pointers); l != 2 {
		t.Errorf("wrong number of pointers, got %d", l)
	}
}
