package geo

import "testing"

func TestPathResample(t *testing.T) {
	p := NewPath()
	p.Resample(10) // should not panic

	p.Push(NewPoint(0, 0)).Resample(10) // should not panic
	p.Push(NewPoint(1.5, 1.5))
	p.Push(NewPoint(2, 2))

	// resample to 0?
	result := p.Clone().Resample(0)
	if result.Length() != 0 {
		t.Error("path, resample down to zero should be empty line")
	}

	// resample to 1
	result = p.Clone().Resample(1)
	answer := NewPath().Push(NewPoint(0, 0))
	if !result.Equals(answer) {
		t.Error("path, resample down to 1 should be first point")
	}

	result = p.Clone().Resample(2)
	answer = NewPath().Push(NewPoint(0, 0)).Push(NewPoint(2, 2))
	if !result.Equals(answer) {
		t.Error("path, resample downsampling")
	}

	result = p.Clone().Resample(5)
	answer = NewPath()
	answer.Push(NewPoint(0, 0)).Push(NewPoint(0.5, 0.5))
	answer.Push(NewPoint(1, 1)).Push(NewPoint(1.5, 1.5))
	answer.Push(NewPoint(2, 2))
	if !result.Equals(answer) {
		t.Error("path, resample upsampling")
		t.Error(result)
		t.Error(answer)
	}

	// round off error case, triggered on my laptop
	p1 := NewPath().Push(NewPoint(-88.145243, 42.321059)).Push(NewPoint(-88.145232, 42.325902))
	p1.Resample(109)
	if p1.Length() != 109 {
		t.Errorf("path, resample incorrect length, expected 109, got %d", p1.Length())
	}

	// duplicate points
	p = NewPath()
	p.Push(NewPoint(1, 0))
	p.Push(NewPoint(1, 0))
	p.Push(NewPoint(1, 0))

	p.Resample(10)
	if l := p.Length(); l != 10 {
		t.Errorf("path, resample length incorrect, got %d", l)
	}

	for i := 0; i < p.Length(); i++ {
		if !p.GetAt(i).Equals(NewPoint(1, 0)) {
			t.Errorf("path, resample not correct point, got %v", p.GetAt(i))
		}
	}
}

func TestPathResampleWithInterval(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(0, 10))

	p.ResampleWithInterval(5.0)
	if l := p.Length(); l != 3 {
		t.Errorf("incorrect resample, got %v", l)
	}

	if v := p.GetAt(1); !v.Equals(NewPoint(0, 5.0)) {
		t.Errorf("incorrect point, got %v", v)
	}
}

func TestPathResampleWithGeoInterval(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0, 0))
	p.Push(NewPoint(0, 10))

	d := p.GeoDistance() / 2
	p.ResampleWithGeoInterval(d)
	if l := p.Length(); l != 3 {
		t.Errorf("incorrect resample, got %v", l)
	}

	if v := p.GetAt(1); !v.Equals(NewPoint(0, 5.0)) {
		t.Errorf("incorrect point, got %v", v)
	}
}

func TestPathResampleEdgeCases(t *testing.T) {
	p := NewPath()
	p.Push(NewPoint(0, 0))

	if !p.resampleEdgeCases(10) {
		t.Errorf("should return false")
	}

	// duplicate points
	p.Push(NewPoint(0, 0))
	if !p.resampleEdgeCases(10) {
		t.Errorf("should return true")
	}

	if l := p.Length(); l != 10 {
		t.Errorf("should reset to suggested points, got %v", l)
	}

	p.resampleEdgeCases(5)
	if l := p.Length(); l != 5 {
		t.Errorf("should shorten if necessary, got %v", l)
	}
}
