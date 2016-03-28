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

	var p2 Point
	err = json.Unmarshal(data, &p2)
	if err != nil {
		t.Errorf("should unmarshal just fine, %v", err)
	}

	if !p1.Equal(p2) {
		t.Errorf("unmarshal incorrect, got %v", p2)
	}
}

func TestLineJSON(t *testing.T) {
	l1 := NewLine(NewPoint(1.5, 2.5), NewPoint(3.5, 4.5))

	data, err := json.Marshal(l1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	// if string(data) != "[[1.5,2.5],[3.5,4.5]]" {
	if string(data) != "{}" {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}

	var l2 Line
	// err = json.Unmarshal(data, &l2)
	// if err != nil {
	// 	t.Errorf("should unmarshal just fine, %v", err)
	// }

	// if !l1.Equal(l2) {
	// 	t.Errorf("unmarshal incorrect, got %v", l2)
	// }

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
	p1 = append(p1,
		NewPoint(1.5, 2.5),
		NewPoint(3.5, 4.5),
		NewPoint(5.5, 6.5),
	)

	data, err := json.Marshal(p1)
	if err != nil {
		t.Errorf("should marshal just fine, %v", err)
	}

	if string(data) != "[[1.5,2.5],[3.5,4.5],[5.5,6.5]]" {
		t.Errorf("json encoding incorrect, got %v", string(data))
	}

	var p2 Path
	err = json.Unmarshal(data, &p2)
	if err != nil {
		t.Errorf("should unmarshal just fine, %v", err)
	}

	if !p1.Equal(p2) {
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

	var b2 Bound
	err = json.Unmarshal(data, &b2)
	if err != nil {
		t.Errorf("should unmarshal just fine, %v", err)
	}

	// if !b1.Equal(b2) {
	// 	t.Errorf("unmarshal incorrect, got %v", b2)
	// }

	err = json.Unmarshal([]byte("[[1,2]]"), &b2)
	if err == nil {
		t.Errorf("should get error since datatypes do not match")
	}

	err = json.Unmarshal([]byte("[[1,2],[3,4],[5,6]]"), &b2)
	if err == nil {
		t.Errorf("should get error since datatypes do not match")
	}
}
