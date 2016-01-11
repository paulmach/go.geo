package reducers

import "testing"

func TestRadialBenchmarkData(t *testing.T) {
	type reduceTest struct {
		Threshold float64
		Length    int
	}

	tests := []reduceTest{
		{0.1, 8282},
		{0.5, 2023},
		{1.0, 1043},
		{1.5, 703},
		{2.0, 527},
		{3.0, 350},
		{4.0, 262},
		{5.0, 209},
	}
	path := benchmarkData()
	for i := range tests {
		p := Radial(path, tests[i].Threshold)
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
		Radial(path, 0.1)
	}
}

func BenchmarkRadialIndexMap(b *testing.B) {
	path := benchmarkData()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RadialIndexMap(path, 0.1)
	}
}
