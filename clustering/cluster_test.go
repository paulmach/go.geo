package clustering

import (
	"testing"

	"github.com/paulmach/go.geo"
)

func TestNewCluster(t *testing.T) {
	// zero pointers
	c1 := NewCluster()

	if c1.Centroid == nil {
		t.Errorf("centroid should not be nil")
	}

	if l := len(c1.Pointers); l != 0 {
		t.Errorf("event list should be empty, got %d", l)
	}

	// one pointer
	c1 = NewCluster(&event{Location: geo.NewPoint(1, 0)})

	if c1.Centroid == c1.Pointers[0].CenterPoint() {
		t.Errorf("should make a copy for center point")
	}

	if !c1.Centroid.Equals(geo.NewPoint(1, 0)) {
		t.Errorf("centroid not adjusted correctly, got %v", c1.Centroid)
	}

	if l := len(c1.Pointers); l != 1 {
		t.Errorf("event not added to list, %d events", l)
	}

	// two pointers
	c1 = NewCluster(
		&event{Location: geo.NewPoint(1, 0)},
		&event{Location: geo.NewPoint(2, 1)},
	)

	if !c1.Centroid.Equals(geo.NewPoint(1.5, 0.5)) {
		t.Errorf("centroid not adjusted correctly, got %v", c1.Centroid)
	}

	if l := len(c1.Pointers); l != 2 {
		t.Errorf("event not added to list, %d events", l)
	}

	// three pointers
	c1 = NewCluster(
		&event{Location: geo.NewPoint(1, 0)},
		&event{Location: geo.NewPoint(2, 1)},
		&event{Location: geo.NewPoint(3, 2)},
	)

	if !c1.Centroid.Equals(geo.NewPoint(2.0, 1.0)) {
		t.Errorf("centroid not adjusted correctly, got %v", c1.Centroid)
	}

	if l := len(c1.Pointers); l != 3 {
		t.Errorf("event not added to list, %d events", l)
	}
}

func TestCombineClusters(t *testing.T) {
	c1 := NewCluster(&event{Location: geo.NewPoint(1, 0)})
	c2 := NewCluster(&event{Location: geo.NewPoint(2, 1)})

	c1 = CombineClusters(c1, c2)
	if !c1.Centroid.Equals(geo.NewPoint(1.5, 0.5)) {
		t.Errorf("centroid not adjusted correctly, got %v", c1.Centroid)
	}

	if l := len(c1.Pointers); l != 2 {
		t.Errorf("event not added to list, %d events", l)
	}

	c3 := NewCluster(&event{Location: geo.NewPoint(3, 2)})
	c1 = CombineClusters(c1, c3)
	if !c1.Centroid.Equals(geo.NewPoint(2.0, 1.0)) {
		t.Errorf("centroid not adjusted correctly, got %v", c1.Centroid)
	}

	if l := len(c1.Pointers); l != 3 {
		t.Errorf("event not added to list, %d events", l)
	}
}
