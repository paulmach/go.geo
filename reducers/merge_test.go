package reducers

import (
	"reflect"
	"testing"
)

func TestMergeIndexMaps(t *testing.T) {
	m1 := []int{0, 1, 3, 5, 7}
	m2 := []int{0, 2, 4}

	merged := MergeIndexMaps(m1, m2)
	if len(merged) != len(m2) {
		t.Errorf("mergeIndexMaps result length incorrect, expected %d, got %d", len(m2), len(merged))
	}

	if merged[0] != 0 {
		t.Errorf("mergeIndexMaps, first value should be zero, got %d", merged[0])
	}

	if merged[len(merged)-1] != m1[len(m1)-1] {
		t.Errorf("mergeIndexMaps, last value should be last value of m1, expected %d, got %d", m1[len(m1)-1], merged[len(merged)-1])
	}

	if !reflect.DeepEqual(merged, []int{0, 3, 7}) {
		t.Errorf("mergeIndexMaps result incorrect, got %v", merged)
	}

}
