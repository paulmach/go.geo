package geo

import "testing"

func BenchmarkPointSetCentroid(b *testing.B) {
	ps := &PointSet{}
	ps.Push(&Point{0, 0}).
		Push(&Point{1, 1.5}).
		Push(&Point{2, 0}).
		Push(&Point{3, 1}).
		Push(&Point{3, 5})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.Centroid()
	}
}

func BenchmarkPointSetGeoCentroid(b *testing.B) {
	ps := &PointSet{}
	ps.Push(&Point{-188.1298828125, -33.97980872872456}).
		Push(&Point{-186.1083984375, -38.54816542304658}).
		Push(&Point{-194.8974609375, -46.10370875598026}).
		Push(&Point{-192.1728515625, -47.8721439688873}).
		Push(&Point{-179.7802734375, -37.30027528134431})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.GeoCentroid()
	}
}
