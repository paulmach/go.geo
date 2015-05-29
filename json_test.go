package geo

import (
	"encoding/json"
	"testing"
)

func TestPointJSON(t *testing.T) {
	p1 := NewPoint(1, 2.1)

	data, err := json.Marshal(p1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	if string(data) != "[1,2.1]" {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}

	var p2 *Point
	err = json.Unmarshal(data, &p2)
	if err != nil {
		t.Errorf("should unmarshal just fine, %v", err)
	}

	if !p1.Equals(p2) {
		t.Errorf("unmarshal incorrect, got %v", p2)
	}
}

func TestLineJSON(t *testing.T) {
	l1 := NewLine(NewPoint(1.5, 2.5), NewPoint(3.5, 4.5))

	data, err := json.Marshal(l1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	if string(data) != "[[1.5,2.5],[3.5,4.5]]" {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}

	var l2 *Line
	err = json.Unmarshal(data, &l2)
	if err != nil {
		t.Errorf("should unmarshal just fine, %v", err)
	}

	if !l1.Equals(l2) {
		t.Errorf("unmarshal incorrect, got %v", l2)
	}

	// decode incomplete data
	err = json.Unmarshal([]byte("[1,2]"), &l2)
	if err == nil {
		t.Errorf("should get error since datatypes don't match")
	}

	err = json.Unmarshal([]byte("[[1,2]]"), &l2)
	if err == nil {
		t.Errorf("should get error since datatypes do not match")
	}

	err = json.Unmarshal([]byte("[[1,2],[3,4],[5,6]]"), &l2)
	if err == nil {
		t.Errorf("should get error since datatypes do not match")
	}
}

func TestPathJSON(t *testing.T) {
	p1 := NewPath()
	p1.Push(NewPoint(1.5, 2.5))
	p1.Push(NewPoint(3.5, 4.5))
	p1.Push(NewPoint(5.5, 6.5))

	data, err := json.Marshal(p1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	if string(data) != "[[1.5,2.5],[3.5,4.5],[5.5,6.5]]" {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}

	var p2 *Path
	err = json.Unmarshal(data, &p2)
	if err != nil {
		t.Errorf("should unmarshal just fine, %v", err)
	}

	if !p1.Equals(p2) {
		t.Errorf("unmarshal incorrect, got %v", p2)
	}

	// empty path
	p1 = NewPath()
	data, err = json.Marshal(p1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	if string(data) != "[]" {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}
}

func TestBoundJSON(t *testing.T) {
	b1 := NewBound(1, 2, 3, 4)

	data, err := json.Marshal(b1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	if string(data) != "[[1,3],[2,4]]" {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}

	var b2 *Bound
	err = json.Unmarshal(data, &b2)
	if err != nil {
		t.Errorf("should unmarshal just fine, %v", err)
	}

	if !b1.Equals(b2) {
		t.Errorf("unmarshal incorrect, got %v", b2)
	}

	err = json.Unmarshal([]byte("[[1,2]]"), &b2)
	if err == nil {
		t.Errorf("should get error since datatypes do not match")
	}

	err = json.Unmarshal([]byte("[[1,2],[3,4],[5,6]]"), &b2)
	if err == nil {
		t.Errorf("should get error since datatypes do not match")
	}
}

func TestSurfaceJSON(t *testing.T) {
	s1 := NewSurface(NewBound(1, 2, 3, 4), 3, 3)
	s1.Grid[0] = []float64{1, 2, 3}
	s1.Grid[1] = []float64{4, 5, 6}
	s1.Grid[2] = []float64{7, 8, 9}

	data, err := json.Marshal(s1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	if string(data) != `{"bound":[[1,3],[2,4]],"values":[[1,2,3],[4,5,6],[7,8,9]]}` {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}

	var s2 *Surface
	err = json.Unmarshal(data, &s2)
	if err != nil {
		t.Errorf("should unmarshal just fine, %v", err)
	}

	if !s1.Bound().Equals(s2.Bound()) {
		t.Errorf("unmarshal of bound incorrect, got %v", s2.Bound())
	}

	if len(s2.Grid) != 3 {
		t.Fatalf("grid not of correct width")
	}

	for i := 0; i < 3; i++ {
		if len(s2.Grid[i]) != 3 {
			t.Fatalf("grid not of correct height")
		}
		for j := 0; j < 3; j++ {
			if s1.Grid[i][j] != s2.Grid[i][j] {
				t.Errorf("grid values incorrect, expected %f, got %f", s1.Grid[i][j], s2.Grid[i][j])
			}
		}
	}
}
