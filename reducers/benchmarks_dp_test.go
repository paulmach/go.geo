package reducers_test

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"testing"

	reducers "."
	"github.com/paulmach/go.geo"
)

func TestDouglasPeuckerBenchmarkData(t *testing.T) {
	type reduceTest struct {
		Threshold float64
		Length    int
	}

	tests := []reduceTest{
		reduceTest{0.1, 1118},
		reduceTest{0.5, 257},
		reduceTest{1.0, 144},
		reduceTest{1.5, 95},
		reduceTest{2.0, 71},
		reduceTest{3.0, 46},
		reduceTest{4.0, 39},
		reduceTest{5.0, 33},
	}
	path := benchmarkData()
	for i := range tests {
		p := reducers.DouglasPeucker(path, tests[i].Threshold)
		if p.Length() != tests[i].Length {
			t.Errorf("douglas peucker benchmark data reduced poorly, got %d, expected %d", p.Length(), tests[i].Length)
		}
	}
}

func BenchmarkDouglasPeucker(b *testing.B) {
	path := benchmarkData()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reducers.DouglasPeucker(path, 0.1)
	}
}

func BenchmarkDouglasPeuckerIndexMap(b *testing.B) {
	path := benchmarkData()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reducers.DouglasPeuckerIndexMap(path, 0.1)
	}
}

func benchmarkData() *geo.Path {
	// Data taken from the simplify-js example at http://mourner.github.io/simplify-js/
	f, err := os.Open("lisbon2portugal.json.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// decompress and decode the json
	var points []float64
	gzipReader, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer gzipReader.Close()

	json.NewDecoder(gzipReader).Decode(&points)

	// create the geo path
	path := geo.NewPathPreallocate(0, len(points)/2)
	for i := 0; i < len(points); i += 2 {
		path.Push(geo.NewPoint(points[i], points[i+1]))
	}

	return path
}
