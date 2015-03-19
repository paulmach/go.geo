package geo

import "testing"

func BenchmarkPointScan(b *testing.B) {
	p := NewPoint(0, 0)
	data := []uint8{1, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64}
	err := p.Scan(data)
	if err != nil {
		b.Fatalf("should scan without error, got %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Scan(data)
	}
}

func BenchmarkPointUnmarshalWKB(b *testing.B) {
	p := NewPoint(0, 0)
	data := []uint8{1, 1, 0, 0, 0, 15, 152, 60, 227, 24, 157, 94, 192, 205, 11, 17, 39, 128, 222, 66, 64}
	err := p.unmarshalWKB(data)
	if err != nil {
		b.Fatalf("should scan without error, got %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.unmarshalWKB(data)
	}
}

func BenchmarkLineScan(b *testing.B) {
	l := &Line{}
	data := []byte{1, 2, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64}

	err := l.Scan(data)
	if err != nil {
		b.Fatalf("should scan without error, got %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Scan(data)
	}
}

func BenchmarkLineUnmarshalWKB(b *testing.B) {
	l := &Line{}
	data := []byte{1, 2, 0, 0, 0, 2, 0, 0, 0, 213, 7, 146, 119, 14, 193, 94, 192, 93, 250, 151, 164, 50, 5, 67, 64, 26, 164, 224, 41, 228, 170, 94, 192, 22, 75, 145, 124, 37, 70, 67, 64}

	err := l.unmarshalWKB(data)
	if err != nil {
		b.Fatalf("should scan without error, got %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.unmarshalWKB(data)
	}
}

func BenchmarkPathScan(b *testing.B) {
	p := NewPath()

	err := p.Scan(testPathWKB)
	if err != nil {
		b.Fatalf("should scan without error, got %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Scan(testPathWKB)
	}
}

func BenchmarkPathUnmarshalWKB(b *testing.B) {
	p := NewPath()

	err := p.unmarshalWKB(testPathWKB)
	if err != nil {
		b.Fatalf("should scan without error, got %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.unmarshalWKB(testPathWKB)
	}
}
