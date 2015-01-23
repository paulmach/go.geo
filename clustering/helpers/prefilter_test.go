package helpers

import (
	"bufio"
	"compress/gzip"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/clustering"
)

func TestRemoveOutlierPointersByQuadkey(t *testing.T) {
	pointers := loadTestPointers(t)

	clusters := RemoveOutlierPointersByQuadkey(pointers, 24, 3)
	if l := len(clusters); l != 555 {
		t.Errorf("incorrect number of clusters, got %v", l)
	}
}

// > go test -c && ./helpers.test -test.bench=RemoveOutlierPointersByQuadkey -test.cpuprofile=cpu.out -test.benchtime=10s
// > go tool pprof prefilter.test cpu.out
func BenchmarkRemoveOutlierPointersByQuadkey(b *testing.B) {
	pointers := loadTestPointers(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RemoveOutlierPointersByQuadkey(pointers, 24, 3)
	}
}

func BenchmarkPrefilteredClusterClustering(b *testing.B) {
	pointers := loadTestPointers(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clusters := RemoveOutlierPointersByQuadkey(pointers, 24, 3)
		clusters = clustering.ClusterClusters(clusters, clustering.CentroidGeoDistance{}, 30)
		if l := len(clusters); l != 27 {
			b.Errorf("incorrect number of clusters, got %v", l)
		}
	}
}

func BenchmarkPrefilteredGeoProjectedClustering(b *testing.B) {
	pointers := loadTestPointers(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clusters := RemoveOutlierPointersByQuadkey(pointers, 24, 3)
		clusters = clustering.ClusterGeoClusters(clusters, 30)

		if l := len(clusters); l != 27 {
			b.Errorf("incorrect number of clusters, got %v", l)
		}
	}
}

func loadTestPointers(tb testing.TB) []clustering.Pointer {
	f, err := os.Open("../testdata/points.csv.gz")
	if err != nil {
		tb.Fatalf("unable to open test file %v", err)
	}
	defer f.Close()

	gzReader, err := gzip.NewReader(f)
	if err != nil {
		tb.Fatalf("unable to create gz reader: %v", err)
	}
	defer gzReader.Close()

	// read in events
	var pointers []clustering.Pointer
	scanner := bufio.NewScanner(gzReader)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		lat, _ := strconv.ParseFloat(parts[0], 64)
		lng, _ := strconv.ParseFloat(parts[1], 64)

		if lat == 0 || lng == 0 {
			tb.Errorf("latlng not parsed correctly, %s %s", parts[0], parts[1])
		}

		pointers = append(pointers, &event{
			Location: geo.NewPoint(lng, lat),
		})
	}

	return pointers
}
