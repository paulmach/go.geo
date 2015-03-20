package clustering

import (
	"testing"

	"github.com/paulmach/go.geo"
)

func TestSortable(t *testing.T) {
	testSet := []*Cluster{
		NewCluster(&event{Location: geo.NewPoint(1, 1)}),
		NewCluster(&event{Location: geo.NewPoint(1, 1)}, &event{Location: geo.NewPoint(2, 2)}, &event{Location: geo.NewPoint(3, 3)}),
		NewCluster(&event{Location: geo.NewPoint(1, 1)}, &event{Location: geo.NewPoint(2, 2)}),
	}

	Sortable(testSet).Sort()

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
