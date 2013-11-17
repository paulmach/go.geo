package geo

import (
	"math"
)

type Projecter struct {
	Project func(p *Point) *Point
	Inverse func(p *Point) *Point
}

var Mercator = Projecter{
	Project: func(p *Point) *Point {
		p.SetLng(math.Pi / 180.0 * p.Lng())

		radLat := math.Pi / 180.0 * p.Lat()
		p.SetLat(math.Log(math.Tan(math.Pi/4 + radLat/2.0)))

		return p
	},
	Inverse: func(p *Point) *Point {
		p.SetLng(p.Lng() * 180.0 / math.Pi)

		radLat := 2.0*math.Atan(math.Exp(p.Lat())) - (math.Pi / 2.0)
		p.SetLat(radLat * 180.0 / math.Pi)

		return p
	},
}

var ScalarMercator struct {
	Level   uint
	Project func(lat, lng float64) (x, y uint)
	Inverse func(x, y uint) (lat, lng float64)
}

func init() {
	ScalarMercator.Level = 31
	ScalarMercator.Project = func(lat, lng float64) (x, y uint) {
		var factor uint
		factor = 1 << ScalarMercator.Level
		maxtiles := float64(factor)

		lng = lng/360.0 + 0.5
		x = (uint)(lng * maxtiles)

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
			y = (uint)(lat * maxtiles)
		}

		return
	}

	ScalarMercator.Inverse = func(x, y uint) (lat, lng float64) {
		factor := 1 << ScalarMercator.Level
		maxtiles := float64(factor)

		lng = 360.0 * (float64(x)/maxtiles - 0.5)
		lat = (2.0*math.Atan(math.Exp(math.Pi-(2*math.Pi)*(float64(y))/maxtiles)))*(180.0/math.Pi) - 90.0

		return
	}
}
