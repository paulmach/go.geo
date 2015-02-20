package geo_test

import (
	"testing"

	geo "."
)

func BenchmarkPointDistanceFrom(b *testing.B) {
	p1 := geo.NewPoint(-122.4167, 37.7833)
	p2 := geo.NewPoint(37.7833, -122.4167)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p1.DistanceFrom(p2)
	}
}

func BenchmarkPointSquaredDistanceFrom(b *testing.B) {
	p1 := geo.NewPoint(-122.4167, 37.7833)
	p2 := geo.NewPoint(37.7833, -122.4167)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p1.SquaredDistanceFrom(p2)
	}
}

func BenchmarkPointQuadKey(b *testing.B) {
	p := geo.NewPoint(-122.4167, 37.7833)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Quadkey(60)
	}
}

func BenchmarkPointQuadKeyString(b *testing.B) {
	p := geo.NewPoint(-122.4167, 37.7833)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.QuadkeyString(60)
	}
}

func BenchmarkPointGeoHash(b *testing.B) {
	p := geo.NewPoint(-122.4167, 37.7833)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.GeoHash()
	}
}

func BenchmarkPointGeoHashInt64(b *testing.B) {
	p := geo.NewPoint(-122.4167, 37.7833)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.GeoHashInt64(60)
	}
}

func BenchmarkPointNormalize(b *testing.B) {
	p := geo.NewPoint(5, 6)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Normalize()
	}
}

func BenchmarkPointEquals(b *testing.B) {
	p1 := geo.NewPoint(5, 6)
	p2 := geo.NewPoint(5, 7)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p1.Equals(p2)
	}
}

func BenchmarkPointClone(b *testing.B) {
	p := geo.NewPoint(5, 6)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Clone()
	}
}
