package geo

import (
	"testing"
)

func TestPointNew(t *testing.T) {
	p := NewPoint(1, 2)
	if p.X() != 1 || p.Lng() != 1 {
		t.Errorf("point, expected 1, got %d", p.X())
	}

	if p.Y() != 2 || p.Lat() != 2 {
		t.Errorf("point, expected 2, got %d", p.Y())
	}
}

func TestPointDistanceFrom(t *testing.T) {
	p1 := NewPoint(0, 0)
	p2 := NewPoint(3, 4)

	if d := p1.DistanceFrom(p2); d != 5 {
		t.Errorf("point, distanceFrom expected 5, got %f", d)
	}

	if d := p2.DistanceFrom(p1); d != 5 {
		t.Errorf("point, distanceFrom expected 5, got %f", d)
	}
}

func TestPointGeoDistanceFrom(t *testing.T) {
	// TODO: implement this test
}

func TestPointBearingTo(t *testing.T) {
	// TODO: implement this test
}

func TestPointAddSubtract(t *testing.T) {
	var answer *Point
	p1 := NewPoint(1, 2)
	p2 := NewPoint(3, 4)

	answer = NewPoint(4, 6)
	if p := p1.Clone().Add(p2); !p.Equals(answer) {
		t.Errorf("point, add expect %v == %v", p, answer)
	}

	answer = NewPoint(-2, -2)
	if p := p1.Clone().Subtract(p2); !p.Equals(answer) {
		t.Errorf("point, subtract expect %v == %v", p, answer)
	}
}

func TestPointNormalize(t *testing.T) {
	var p, answer *Point

	p = NewPoint(5, 0)
	answer = NewPoint(1, 0)
	if p.Normalize(); !p.Equals(answer) {
		t.Errorf("point, normalize expect %v == %v", p, answer)
	}

	p = NewPoint(0, 5)
	answer = NewPoint(0, 1)
	if p.Normalize(); !p.Equals(answer) {
		t.Errorf("point, normalize expect %v == %v", p, answer)
	}
}

func TestPointScale(t *testing.T) {
	var p, answer *Point

	p = NewPoint(5, 0)
	answer = NewPoint(10, 0)
	if p.Scale(2.0); !p.Equals(answer) {
		t.Errorf("point, scale expect %v == %v", p, answer)
	}

	p = NewPoint(0, 5)
	answer = NewPoint(0, 15)
	if p.Scale(3.0); !p.Equals(answer) {
		t.Errorf("point, scale expect %v == %v", p, answer)
	}

	p = NewPoint(2, 3)
	answer = NewPoint(10, 15)
	if p.Scale(5.0); !p.Equals(answer) {
		t.Errorf("point, scale expect %v == %v", p, answer)
	}

	p = NewPoint(2, 3)
	answer = NewPoint(-10, -15)
	if p.Scale(-5.0); !p.Equals(answer) {
		t.Errorf("point, scale expect %v == %v", p, answer)
	}
}

func TestDot(t *testing.T) {
	p1 := NewPoint(0, 0)

	p2 := NewPoint(1, 2)
	answer := 0.0
	if d := p1.Dot(p2); d != answer {
		t.Errorf("point, dot expteced %v == %v", d, answer)
	}

	p1 = NewPoint(4, 5)
	answer = 14.0
	if d := p1.Dot(p2); d != answer {
		t.Errorf("point, dot expteced %v == %v", d, answer)
	}

	// reverse version
	if d := p2.Dot(p1); d != answer {
		t.Errorf("point, dot expteced %v == %v", d, answer)
	}
}

func TestPointClone(t *testing.T) {
	p1 := NewPoint(1, 0)
	p2 := NewPoint(1, 2)

	if p := p1.Clone(); !p.Equals(p1) {
		t.Errorf("point, clone expect %v == %v", p, p1)
	}

	if p := p1.Clone(); p.Equals(p2) {
		t.Errorf("point, clone expect %v != %v", p, p2)
	}

	if p := p1.Clone().SetX(10); p.Equals(p1) {
		t.Errorf("point, clone expect %v != %v", p, p1)
	}
}

func TestPointEquals(t *testing.T) {
	p1 := NewPoint(1, 0)
	p2 := NewPoint(1, 0)

	p3 := NewPoint(2, 3)
	p4 := NewPoint(2, 4)

	if !p1.Equals(p2) {
		t.Errorf("point, equals expect %v == %v", p1, p2)
	}

	if p2.Equals(p3) {
		t.Errorf("point, equals expect %v != %v", p2, p3)
	}

	if p3.Equals(p4) {
		t.Errorf("point, equals expect %v != %v", p3, p4)
	}
}

func TestPointGettersSetters(t *testing.T) {
	var p *Point

	p = NewPoint(0, 0)
	p.SetX(10)
	if v := p.X(); v != 10 {
		t.Errorf("point, setX expect %f == 10", v)
	}

	if p.Y() != 0 {
		t.Error("point, setX expected Y to be unchanged")
	}

	p = NewPoint(1, 0)
	p.SetY(5)
	if v := p.Y(); v != 5 {
		t.Errorf("point, setY expect %f == 5", v)
	}

	if p.X() != 1 {
		t.Error("point, setY expected X to be unchanged")
	}

	p = NewPoint(2, 0)
	p.SetLat(-12.3)
	if v := p.Lat(); v != -12.3 {
		t.Errorf("point, setLat expect %f == -12.3", v)
	}

	if p.Lng() != 2 {
		t.Error("point, setLat expected Lng to be unchanged")
	}

	p = NewPoint(0, 3)
	p.SetLng(-45.6)
	if v := p.Lng(); v != -45.6 {
		t.Errorf("point, setLng expect %f == -45.6", v)
	}

	if p.Lat() != 3 {
		t.Error("point, setLng expected Lat to be unchanged")
	}
}

func TestPointToArray(t *testing.T) {
	p := NewPoint(1, 2)
	if a := p.ToArray(); a != [2]float64{1, 2} {
		t.Errorf("point, toArray expected %v == %v", a, [2]float64{1, 2})
	}
}

func TestPointString(t *testing.T) {
	p := NewPoint(1, 2)

	answer := "[1.000000, 2.000000]"
	if s := p.String(); s != answer {
		t.Errorf("point, string expected %s, got %s", answer, s)
	}
}
