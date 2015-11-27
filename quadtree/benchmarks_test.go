package quadtree

import (
	"math"
	"math/rand"
	"testing"

	"github.com/paulmach/go.geo"
)

func BenchmarkInsert(b *testing.B) {
	r := rand.New(rand.NewSource(22))
	qt := New(geo.NewBound(0, 1, 0, 1))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qt.Insert(geo.NewPoint(r.Float64(), r.Float64()))
	}
}

func BenchmarkFromPointer50(b *testing.B) {
	r := rand.New(rand.NewSource(32))
	pointers := make([]geo.Pointer, 0, 50)
	for i := 0; i < 50; i++ {
		pointers = append(pointers, geo.NewPoint(r.Float64(), r.Float64()))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFromPointers(pointers)
	}
}

func BenchmarkFromPointer100(b *testing.B) {
	r := rand.New(rand.NewSource(42))
	pointers := make([]geo.Pointer, 0, 100)
	for i := 0; i < 100; i++ {
		pointers = append(pointers, geo.NewPoint(r.Float64(), r.Float64()))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFromPointers(pointers)
	}
}

func BenchmarkFromPointer500(b *testing.B) {
	r := rand.New(rand.NewSource(52))
	pointers := make([]geo.Pointer, 0, 500)
	for i := 0; i < 500; i++ {
		pointers = append(pointers, geo.NewPoint(r.Float64(), r.Float64()))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFromPointers(pointers)
	}
}

func BenchmarkFromPointer1000(b *testing.B) {
	r := rand.New(rand.NewSource(62))
	pointers := make([]geo.Pointer, 0, 1000)
	for i := 0; i < 1000; i++ {
		pointers = append(pointers, geo.NewPoint(r.Float64(), r.Float64()))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFromPointers(pointers)
	}
}

func BenchmarkRandomFind1000(b *testing.B) {
	r := rand.New(rand.NewSource(42))

	var pointers []geo.Pointer
	for i := 0; i < 1000; i++ {
		pointers = append(pointers, geo.NewPoint(r.Float64(), r.Float64()))
	}

	qt := NewFromPointers(pointers)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qt.Find(geo.NewPoint(r.Float64(), r.Float64()))
	}
}

func BenchmarkRandomFind1000Naive(b *testing.B) {
	r := rand.New(rand.NewSource(42))

	var pointers []geo.Pointer
	for i := 0; i < 1000; i++ {
		pointers = append(pointers, geo.NewPoint(r.Float64(), r.Float64()))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		looking := geo.NewPoint(r.Float64(), r.Float64())

		min := math.MaxFloat64
		var best geo.Pointer
		for _, p := range pointers {
			if d := looking.SquaredDistanceFrom(p.Point()); d < min {
				min = d
				best = p
			}
		}

		_ = best
	}
}

func BenchmarkRandomInBound1000(b *testing.B) {
	r := rand.New(rand.NewSource(43))

	var pointers []geo.Pointer
	for i := 0; i < 1000; i++ {
		pointers = append(pointers, geo.NewPoint(r.Float64(), r.Float64()))
	}

	qt := NewFromPointers(pointers)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := geo.NewPoint(r.Float64(), r.Float64())
		qt.InBound(geo.NewBoundFromPoints(p, p).Pad(0.1))
	}
}

func BenchmarkRandomInBound1000Naive(b *testing.B) {
	r := rand.New(rand.NewSource(43))

	var pointers []geo.Pointer
	for i := 0; i < 1000; i++ {
		pointers = append(pointers, geo.NewPoint(r.Float64(), r.Float64()))
	}

	b.ReportAllocs()
	b.ResetTimer()

	var near []geo.Pointer
	for i := 0; i < b.N; i++ {
		p := geo.NewPoint(r.Float64(), r.Float64())
		b := geo.NewBoundFromPoints(p, p).Pad(0.1)

		near = near[:0]
		for _, p := range pointers {
			if b.Contains(p.Point()) {
				near = append(near, p)
			}
		}

		_ = len(near)
	}
}

func BenchmarkRandomInBound1000Buf(b *testing.B) {
	r := rand.New(rand.NewSource(43))

	var pointers []geo.Pointer
	for i := 0; i < 1000; i++ {
		pointers = append(pointers, geo.NewPoint(r.Float64(), r.Float64()))
	}

	qt := NewFromPointers(pointers)

	var buf []geo.Pointer
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := geo.NewPoint(r.Float64(), r.Float64())
		buf = qt.InBound(geo.NewBoundFromPoints(p, p).Pad(0.1), buf)
	}
}
