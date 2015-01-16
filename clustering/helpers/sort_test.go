package helpers

import (
	"testing"

	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/clustering"
)

func TestSortableClusters(t *testing.T) {
	testSet := []*clustering.Cluster{
		clustering.NewCluster(&event{Location: geo.NewPoint(1, 1)}),
		clustering.NewCluster(&event{Location: geo.NewPoint(1, 1)}, &event{Location: geo.NewPoint(2, 2)}, &event{Location: geo.NewPoint(3, 3)}),
		clustering.NewCluster(&event{Location: geo.NewPoint(1, 1)}, &event{Location: geo.NewPoint(2, 2)}),
	}

	SortableClusters(testSet).Sort()

	if l := len(testSet); l != 3 {
		t.Errorf("set length incorrect, got %d", l)
	}

	if l := len(testSet[0].Pointers); l != 3 {
		t.Errorf("length of set incorrect, got %d", l)
	}

	if l := len(testSet[1].Pointers); l != 2 {
		t.Errorf("length of set incorrect, got %d", l)
	}

	if l := len(testSet[2].Pointers); l != 1 {
		t.Errorf("length of set incorrect, got %d", l)
	}
}
