package geo_test

import (
	"testing"

	geo "."
)

func BenchmarkLineDistanceFrom(b *testing.B) {
	l := geo.NewLine(geo.NewPoint(1, 2), geo.NewPoint(3, 4))
	p := geo.NewPoint(2, 4)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.DistanceFrom(p)
	}
}

func BenchmarkLineSquaredDistanceFrom(b *testing.B) {
	l := geo.NewLine(geo.NewPoint(1, 2), geo.NewPoint(3, 4))
	p := geo.NewPoint(2, 4)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.SquaredDistanceFrom(p)
	}
}

func BenchmarkLineProject(b *testing.B) {
	l := geo.NewLine(geo.NewPoint(1, 2), geo.NewPoint(3, 4))
	p := geo.NewPoint(2, 4)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Project(p)
	}
}

func BenchmarkLineMeasure(b *testing.B) {
	l := geo.NewLine(geo.NewPoint(1, 2), geo.NewPoint(3, 4))
	p := geo.NewPoint(5, 4)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Measure(p)
	}
}

func BenchmarkLineInterpolate(b *testing.B) {
	l := geo.NewLine(geo.NewPoint(1, 2), geo.NewPoint(3, 4))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Interpolate(0.5)
	}
}

func BenchmarkLineMidpoint(b *testing.B) {
	l := geo.NewLine(geo.NewPoint(1, 2), geo.NewPoint(3, 4))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Midpoint()
	}
}

func BenchmarkLineGeoMidpoint(b *testing.B) {
	l := geo.NewLine(geo.NewPoint(1, 2), geo.NewPoint(3, 4))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.GeoMidpoint()
	}
}

func BenchmarkLineEquals(b *testing.B) {
	l1 := geo.NewLine(geo.NewPoint(1, 2), geo.NewPoint(3, 4))
	l2 := geo.NewLine(geo.NewPoint(1, 2), geo.NewPoint(4, 3))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l1.Equals(l2)
	}
}

func BenchmarkLineClone(b *testing.B) {
	l := geo.NewLine(geo.NewPoint(1, 2), geo.NewPoint(3, 4))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Clone()
	}
}
