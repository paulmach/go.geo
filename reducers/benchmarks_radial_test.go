package reducers_test

import (
	"testing"

	reducers "."
)

func TestRadialBenchmarkData(t *testing.T) {
	type reduceTest struct {
		Threshold float64
		Length    int
	}

	tests := []reduceTest{
		reduceTest{0.1, 8282},
		reduceTest{0.5, 2023},
		reduceTest{1.0, 1043},
		reduceTest{1.5, 703},
		reduceTest{2.0, 527},
		reduceTest{3.0, 350},
		reduceTest{4.0, 262},
		reduceTest{5.0, 209},
	}
	path := benchmarkData()
	for i := range tests {
		p := reducers.Radial(path, tests[i].Threshold)
		if p.Length() != tests[i].Length {
			t.Errorf("radial benchmark data reduced poorly, got %d, expected %d", p.Length(), tests[i].Length)
		}
	}
}

func BenchmarkRadial(b *testing.B) {
	path := benchmarkData()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reducers.Radial(path, 0.1)
	}
}

func BenchmarkRadialIndexMap(b *testing.B) {
	path := benchmarkData()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reducers.RadialIndexMap(path, 0.1)
	}
}
