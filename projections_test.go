package geo

import (
	"math"
	"testing"
)

var cities = [][2]float64{
	{57.09700, 9.85000}, {49.03000, -122.32000}, {39.23500, -76.17490},
	{57.20000, -2.20000}, {16.75000, -99.76700}, {5.60000, -0.16700},
	{51.66700, -176.46700}, {9.00000, 38.73330}, {-34.7666, 138.53670},
	{12.80000, 45.00000}, {42.70000, -110.86700}, {13.48167, 144.79330},
	{33.53300, -81.71700}, {42.53300, -99.85000}, {26.01670, 50.55000},
	{35.75000, -84.00000}, {51.11933, -1.15543}, {82.52000, -62.28000},
	{32.91700, -85.91700}, {31.19000, 29.95000}, {36.70000, 3.21700},
	{34.14000, -118.10700}, {32.50370, -116.45100}, {47.83400, 10.86800},
	{28.25000, 129.70000}, {16.75000, -22.95000}, {31.95000, 35.95000},
	{52.35000, 4.86660}, {13.58670, 144.93670}, {6.90000, 134.15000},
	{40.03000, 32.90000}, {33.65000, -85.78300}, {49.33000, 10.59700},
	{17.13330, -61.78330}, {-23.4333, -70.60000}, {51.21670, 4.40000},
	{29.60000, 35.01000}, {38.58330, -121.48300}, {34.16700, -97.13300},
	{45.60000, 9.15000}, {-18.3500, -70.33330}, {-7.88000, -14.42000},
	{15.28330, 38.90000}, {-25.2333, -57.51670}, {23.96500, 32.82000},
	{-36.8832, 174.75000}, {-38.0333, 144.46670}, {46.03300, 12.60000},
	{41.66700, -72.83300}, {35.45000, 139.45000}}

func TestMercator(t *testing.T) {
	for _, city := range cities {
		p := &Point{}

		p.SetLat(city[0])
		p.SetLng(city[1])

		Mercator.Project(p)
		Mercator.Inverse(p)

		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("Mercator, latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lng()-city[1]) > epsilon {
			t.Errorf("Mercator, longitude miss match: %f != %f", p.Lng(), city[1])
		}
	}
}

func TestTransverseMercator(t *testing.T) {
	tested := 0

	for _, city := range cities {
		p := &Point{}

		p.SetLat(city[0])
		p.SetLng(city[1])

		if math.Abs(p.Lng()) > 10 {
			continue
		}

		TransverseMercator.Project(p)
		TransverseMercator.Inverse(p)

		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("TransverseMercator, latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lng()-city[1]) > epsilon {
			t.Errorf("TransverseMercator, longitude miss match: %f != %f", p.Lng(), city[1])
		}

		tested++
	}

	if tested == 0 {
		t.Error("TransverseMercator, no points tested")
	}
}

func TestBuildTransverseMercator(t *testing.T) {
	for _, city := range cities {
		p := &Point{}

		p.SetLat(city[0])
		p.SetLng(city[1])

		offset := math.Floor(p.Lng()/10.0) * 10.0
		projector := BuildTransverseMercator(offset)

		projector.Project(p)
		projector.Inverse(p)

		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("BuildTransverseMercator, latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lng()-city[1]) > epsilon {
			t.Errorf("BuildTransverseMercator, longitude miss match: %f != %f", p.Lng(), city[1])
		}
	}

	// test anti-meridian from right
	projector := BuildTransverseMercator(-178.0)

	test := NewPoint(-175.0, 30)

	p := test.Clone()
	projector.Project(p)
	projector.Inverse(p)

	if math.Abs(p.Lat()-test.Lat()) > epsilon {
		t.Errorf("BuildTransverseMercator, latitude miss match: %f != %f", p.Lat(), test.Lat())
	}

	if math.Abs(p.Lng()-test.Lng()) > epsilon {
		t.Errorf("BuildTransverseMercator, longitude miss match: %f != %f", p.Lng(), test.Lat())
	}

	test = NewPoint(179.0, 30)

	p = test.Clone()
	projector.Project(p)
	projector.Inverse(p)

	if math.Abs(p.Lat()-test.Lat()) > epsilon {
		t.Errorf("BuildTransverseMercator, latitude miss match: %f != %f", p.Lat(), test.Lat())
	}

	if math.Abs(p.Lng()-test.Lng()) > epsilon {
		t.Errorf("BuildTransverseMercator, longitude miss match: %f != %f", p.Lng(), test.Lat())
	}

	// test anti-meridian from left
	projector = BuildTransverseMercator(178.0)

	test = NewPoint(175.0, 30)

	p = test.Clone()
	projector.Project(p)
	projector.Inverse(p)

	if math.Abs(p.Lat()-test.Lat()) > epsilon {
		t.Errorf("BuildTransverseMercator, latitude miss match: %f != %f", p.Lat(), test.Lat())
	}

	if math.Abs(p.Lng()-test.Lng()) > epsilon {
		t.Errorf("BuildTransverseMercator, longitude miss match: %f != %f", p.Lng(), test.Lat())
	}

	test = NewPoint(-179.0, 30)

	p = test.Clone()
	projector.Project(p)
	projector.Inverse(p)

	if math.Abs(p.Lat()-test.Lat()) > epsilon {
		t.Errorf("BuildTransverseMercator, latitude miss match: %f != %f", p.Lat(), test.Lat())
	}

	if math.Abs(p.Lng()-test.Lng()) > epsilon {
		t.Errorf("BuildTransverseMercator, longitude miss match: %f != %f", p.Lng(), test.Lat())
	}
}

func TestScalarMercator(t *testing.T) {

	x, y := ScalarMercator.Project(0, 0)
	lat, lng := ScalarMercator.Inverse(x, y)

	if lat != 0.0 {
		t.Errorf("Scalar Mercator, latitude should be 0: %f", lat)
	}

	if lng != 0.0 {
		t.Errorf("Scalar Mercator, longitude should be 0: %f", lng)
	}

	for _, city := range cities {
		p := &Point{}

		p.SetLat(city[0])
		p.SetLng(city[1])

		x, y := ScalarMercator.Project(p.Lng(), p.Lat())
		lng, lat := ScalarMercator.Inverse(x, y)

		p.SetLat(lat)
		p.SetLng(lng)

		if math.Abs(p.Lat()-city[0]) > epsilon {
			t.Errorf("Scalar Mercator, latitude miss match: %f != %f", p.Lat(), city[0])
		}

		if math.Abs(p.Lng()-city[1]) > epsilon {
			t.Errorf("Scalar Mercator, longitude miss match: %f != %f", p.Lng(), city[1])
		}
	}

	// test polar regions
	if _, y := ScalarMercator.Project(0, 89.9); y != (1<<ScalarMercator.Level)-1 {
		t.Error("Scalar Mercator, top of the world error")
	}

	if _, y := ScalarMercator.Project(0, -89.9); y != 0 {
		t.Error("Scalar Mercator, bottom of the world error")
	}
}
