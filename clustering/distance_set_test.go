package clustering

import "testing"

func TestNewDistanceSet(t *testing.T) {
	s := newDistanceSet()

	if s.Distances == nil {
		t.Errorf("must initialize map")
	}

	if len(s.Distances) != 0 {
		t.Errorf("must initialize with empty map")
	}
}

func TestDistanceSetSet(t *testing.T) {
	s := newDistanceSet()

	s.Set(1, 10)
	if v := s.MinIndex; v != 1 {
		t.Errorf("incorrect min index, got %v", v)
	}

	if v := s.MinDistance; v != 10 {
		t.Errorf("incorrect min distance, got %v", v)
	}

	// add a higher distance
	s.Set(2, 20)
	if v := s.MinIndex; v != 1 {
		t.Errorf("incorrect min index, got %v", v)
	}

	if v := s.MinDistance; v != 10 {
		t.Errorf("incorrect min distance, got %v", v)
	}

	// add a lower distance
	s.Set(3, 5)
	if v := s.MinIndex; v != 3 {
		t.Errorf("incorrect min index, got %v", v)
	}

	if v := s.MinDistance; v != 5 {
		t.Errorf("incorrect min distance, got %v", v)
	}

	// set a lower distance for the current min index
	s.Set(3, 2)
	if v := s.MinIndex; v != 3 {
		t.Errorf("incorrect min index, got %v", v)
	}

	if v := s.MinDistance; v != 2 {
		t.Errorf("incorrect min distance, got %v", v)
	}

	// set a very high distance for the current min index
	s.Set(3, 50)
	if v := s.MinIndex; v != 1 {
		t.Errorf("incorrect min index, got %v", v)
	}

	if v := s.MinDistance; v != 10 {
		t.Errorf("incorrect min distance, got %v", v)
	}
}

func TestDistanceSetDelete(t *testing.T) {
	s := newDistanceSet()

	s.Set(1, 10)
	s.Set(2, 20)
	s.Set(3, 5)
	s.Set(4, 30)

	if v := s.MinIndex; v != 3 {
		t.Errorf("incorrect min index, got %v", v)
	}

	if v := s.MinDistance; v != 5 {
		t.Errorf("incorrect min distance, got %v", v)
	}

	// delete not min
	s.Delete(1)

	if v := s.MinIndex; v != 3 {
		t.Errorf("incorrect min index, got %v", v)
	}

	if v := s.MinDistance; v != 5 {
		t.Errorf("incorrect min distance, got %v", v)
	}

	// delete min
	s.Delete(s.MinIndex)

	if v := s.MinIndex; v != 2 {
		t.Errorf("incorrect min index, got %v", v)
	}

	if v := s.MinDistance; v != 20 {
		t.Errorf("incorrect min distance, got %v", v)
	}
}
