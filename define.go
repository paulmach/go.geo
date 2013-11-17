package geo

import (
	"math"
)

var epsilon = 1e-6

var DEFAULT_HaversineGeoDistance = false // or use linear approximation
var DEFAULT_Radius = 6371000.0

func yesHaversine(haversine []bool) bool {
	return (len(haversine) != 0 && haversine[0]) || (DEFAULT_HaversineGeoDistance && len(haversine) == 0)
}

func deg2rad(d float64) float64 {
	return d * math.Pi / 180.0
}

func rad2deg(r float64) float64 {
	return 180.0 * r / math.Pi
}
