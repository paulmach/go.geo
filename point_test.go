package geo

import (
	"math"
	"strings"
	"testing"
)

var citiesGeoHash = [][3]interface{}{
	{57.09700, 9.85000, "u4phb4hw"},
	{49.03000, -122.32000, "c29nbt9k3q"},
	{39.23500, -76.17490, "dqcz4we0k"},
	{-34.7666, 138.53670, "r1fd0qzmg"},
}

func TestNewPoint(t *testing.T) {
	p := NewPoint(1, 2)
	if p.X() != 1 || p.Lng() != 1 {
		t.Errorf("point, expected 1, got %f", p.X())
	}

	if p.Y() != 2 || p.Lat() != 2 {
		t.Errorf("point, expected 2, got %f", p.Y())
	}
}

func TestPointQuadkey(t *testing.T) {
	p := &Point{}

	p.SetLat(41.850033)
	p.SetLng(-87.65005229999997)

	if k := p.Quadkey(15); k != 212521785 {
		t.Errorf("point quadkey, incorrect got %d", k)
	}

	// default level
	level := 30
	for _, city := range cities {
		p := &Point{}

		p.SetLat(city[0])
		p.SetLng(city[1])

		key := p.Quadkey(level)

		p = NewPointFromQuadkey(key, level)

		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("point quadkey, latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lng()-city[1]) > epsilon {
			t.Errorf("point quadkey, longitude miss match: %f != %f", p.Lng(), city[1])
		}
	}
}

func TestPointQuadkeyString(t *testing.T) {
	p := &Point{}

	p.SetLat(41.850033)
	p.SetLng(-87.65005229999997)

	if k := p.QuadkeyString(15); k != "030222231030321" {
		t.Errorf("point quadkey string, incorrect got %s", k)
	}

	// default level
	level := 30
	for _, city := range cities {
		p := &Point{}

		p.SetLat(city[0])
		p.SetLng(city[1])

		key := p.QuadkeyString(level)

		p = NewPointFromQuadkeyString(key)

		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("point quadkey, latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lng()-city[1]) > epsilon {
			t.Errorf("point quadkey, longitude miss match: %f != %f", p.Lng(), city[1])
		}
	}
}

func TestNewPointFromGeoHash(t *testing.T) {
	for _, c := range citiesGeoHash {
		p := NewPointFromGeoHash(c[2].(string))
		if d := p.GeoDistanceFrom(NewPoint(c[1].(float64), c[0].(float64))); d > 10 {
			t.Errorf("point, new from geohash expected distance %f", d)
		}
	}
}

func TestNewPointFromGeoHashInt64(t *testing.T) {
	for _, c := range citiesGeoHash {
		var hash int64
		for _, r := range c[2].(string) {
			hash <<= 5
			hash |= int64(strings.Index("0123456789bcdefghjkmnpqrstuvwxyz", string(r)))
		}

		p := NewPointFromGeoHashInt64(hash, 5*len(c[2].(string)))
		if d := p.GeoDistanceFrom(NewPoint(c[1].(float64), c[0].(float64))); d > 10 {
			t.Errorf("point, new from geohash expected distance %f", d)
		}
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

func TestPointSquaredDistanceFrom(t *testing.T) {
	p1 := NewPoint(0, 0)
	p2 := NewPoint(3, 4)

	if d := p1.SquaredDistanceFrom(p2); d != 25 {
		t.Errorf("point, squaredDistanceFrom expected 25, got %f", d)
	}

	if d := p2.SquaredDistanceFrom(p1); d != 25 {
		t.Errorf("point, squaredDistanceFrom expected 25, got %f", d)
	}
}

func TestPointGeoDistanceFrom(t *testing.T) {
	p1 := NewPoint(-1.8444, 53.1506)
	p2 := NewPoint(0.1406, 52.2047)

	if d := p1.GeoDistanceFrom(p2, true); math.Abs(d-170389.801924) > epsilon {
		t.Errorf("incorrect geodistance, got %v", d)
	}

	if d := p1.GeoDistanceFrom(p2, false); math.Abs(d-170400.503437) > epsilon {
		t.Errorf("incorrect geodistance, got %v", d)
	}

	p1 = NewPoint(0.5, 30)
	p2 = NewPoint(-0.5, 30)

	dFast := p1.GeoDistanceFrom(p2, false)
	dHav := p1.GeoDistanceFrom(p2, true)

	p1 = NewPoint(179.5, 30)
	p2 = NewPoint(-179.5, 30)

	if d := p1.GeoDistanceFrom(p2, false); math.Abs(d-dFast) > epsilon {
		t.Errorf("incorrect geodistance, got %v", d)
	}

	if d := p1.GeoDistanceFrom(p2, true); math.Abs(d-dHav) > epsilon {
		t.Errorf("incorrect geodistance, got %v", d)
	}
}

func TestPointBearingTo(t *testing.T) {
	p1 := NewPoint(0, 0)
	p2 := NewPoint(0, 1)

	if d := p1.BearingTo(p2); d != 0 {
		t.Errorf("point, bearingTo expected 0, got %f", d)
	}

	if d := p2.BearingTo(p1); d != 180 {
		t.Errorf("point, bearingTo expected 180, got %f", d)
	}

	p1 = NewPoint(0, 0)
	p2 = NewPoint(1, 0)

	if d := p1.BearingTo(p2); d != 90 {
		t.Errorf("point, bearingTo expected 90, got %f", d)
	}

	if d := p2.BearingTo(p1); d != -90 {
		t.Errorf("point, bearingTo expected -90, got %f", d)
	}

	p1 = NewPoint(-1.8444, 53.1506)
	p2 = NewPoint(0.1406, 52.2047)

	if d := p1.BearingTo(p2); math.Abs(127.373351-d) > epsilon {
		t.Errorf("point, bearingTo got %f", d)
	}
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

	p = NewPoint(0, 0)
	answer = NewPoint(0, 0)
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

func TestPointDot(t *testing.T) {
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

func TestPointGeoHash(t *testing.T) {
	for _, c := range citiesGeoHash {
		hash := NewPoint(c[1].(float64), c[0].(float64)).GeoHash()
		if !strings.HasPrefix(hash, c[2].(string)) {
			t.Errorf("point, geohash expected %s, got %s", c[2].(string), hash)
		}
	}

	for _, c := range citiesGeoHash {
		hash := NewPoint(c[1].(float64), c[0].(float64)).GeoHash(len(c[2].(string)))
		if hash != c[2].(string) {
			t.Errorf("point, geohash expected %s, got %s", c[2].(string), hash)
		}
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

func TestPointToGeoJSON(t *testing.T) {
	p := NewPoint(1, 2.5)

	f := p.ToGeoJSON()
	if !f.Geometry.IsPoint() {
		t.Errorf("point, should be point geometry")
	}
}

func TestPointToWKT(t *testing.T) {
	p := NewPoint(1, 2.5)

	answer := "POINT(1 2.5)"
	if s := p.ToWKT(); s != answer {
		t.Errorf("point, string expected %s, got %s", answer, s)
	}
}

func TestPointString(t *testing.T) {
	p := NewPoint(1, 2.5)

	answer := "POINT(1 2.5)"
	if s := p.String(); s != answer {
		t.Errorf("point, string expected %s, got %s", answer, s)
	}
}
