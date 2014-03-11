package geo

import (
	"math"
)

type Projecter struct {
	Project func(p *Point)
	Inverse func(p *Point)
}

var Mercator = Projecter{
	Project: func(p *Point) {
		p.SetX(math.Pi / 180.0 * p.Lng())

		radLat := math.Pi / 180.0 * p.Lat()
		p.SetY(math.Log(math.Tan(math.Pi/4 + radLat/2.0)))
	},
	Inverse: func(p *Point) {
		p.SetLng(p.X() * 180.0 / math.Pi)

		radLat := 2.0*math.Atan(math.Exp(p.Y())) - (math.Pi / 2.0)
		p.SetLat(radLat * 180.0 / math.Pi)
	},
}

// BuildTransverseMercator builds a transverse mercator projection
// that automatically recenter the longitude around the provided centerLng.
// Works correctly around the anti-meridian.
// http://en.wikipedia.org/wiki/Transverse_Mercator_projection
func BuildTransverseMercator(centerLng float64) Projecter {
	return Projecter{
		Project: func(p *Point) {
			lng := p.Lng() - centerLng
			if lng < 180 {
				lng += 360.0
			}

			if lng > 180 {
				lng -= 360.0
			}

			p.SetLng(lng)
			TransverseMercator.Project(p)
		},
		Inverse: func(p *Point) {
			TransverseMercator.Inverse(p)

			lng := p.Lng() + centerLng
			if lng < 180 {
				lng += 360.0
			}

			if lng > 180 {
				lng -= 360.0
			}

			p.SetLng(lng)
		},
	}
}

// This default Transverse Mercator projector will only work well +-10 degrees around
// longitude 0. Use this if you've already precentered your points.
var TransverseMercator = Projecter{
	Project: func(p *Point) {
		radLat := deg2rad(p.Lat())
		radLng := deg2rad(p.Lng())

		sincos := math.Sin(radLng) * math.Cos(radLat)
		p.SetX(0.5 * math.Log((1+sincos)/(1-sincos)))

		p.SetY(math.Atan(math.Tan(radLat) / math.Cos(radLng)))
	},
	Inverse: func(p *Point) {
		lng := math.Atan(math.Sinh(p.X()) / math.Cos(p.Y()))
		lat := math.Asin(math.Sin(p.Y()) / math.Cosh(p.X()))

		p.SetLng(rad2deg(lng))
		p.SetLat(rad2deg(lat))
	},
}

// ScalarMercator projects converts from lng/lat float64 to x,y uint64.
// This is similar to Google's world coordinates.
var ScalarMercator struct {
	Level   uint64
	Project func(lat, lng float64, level ...uint64) (x, y uint64)
	Inverse func(x, y uint64, level ...uint64) (lat, lng float64)
}

func init() {
	ScalarMercator.Level = 31
	ScalarMercator.Project = scalarMercatorProject
	ScalarMercator.Inverse = scalarMercatorInverse
}

func scalarMercatorProject(lng, lat float64, level ...uint64) (x, y uint64) {
	var factor uint64
	l := ScalarMercator.Level
	if len(level) != 0 {
		l = level[0]
	}

	factor = 1 << l
	maxtiles := float64(factor)

	lng = lng/360.0 + 0.5
	x = (uint64)(lng * maxtiles)

	// bound it because we have a top of the world problem
	siny := math.Sin(lat * math.Pi / 180.0)

	if siny < -0.9999 {
		lat = 0.5 + 0.5*math.Log((1.0+siny)/(1.0-siny))/(-2*math.Pi)
		y = 0
	} else if siny > 0.9999 {
		lat = 0.5 + 0.5*math.Log((1.0+siny)/(1.0-siny))/(-2*math.Pi)
		y = factor - 1
	} else {
		lat = 0.5 + 0.5*math.Log((1.0+siny)/(1.0-siny))/(-2*math.Pi)
		y = (uint64)(lat * maxtiles)
	}

	return
}

func scalarMercatorInverse(x, y uint64, level ...uint64) (lng, lat float64) {
	var factor uint64
	l := ScalarMercator.Level
	if len(level) != 0 {
		l = level[0]
	}

	factor = 1 << l
	maxtiles := float64(factor)

	lng = 360.0 * (float64(x)/maxtiles - 0.5)
	lat = (2.0*math.Atan(math.Exp(math.Pi-(2*math.Pi)*(float64(y))/maxtiles)))*(180.0/math.Pi) - 90.0

	return
}
