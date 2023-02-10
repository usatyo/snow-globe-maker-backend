package main

import (
	"math"
)

func calcDist(latitude1, longitude1, alt1, latitude2, longitude2, alt2 float64) float64 {
	// 高度を考慮した距離を計算する
	dist := haversine(latitude1, longitude1, latitude2, longitude2)
	altDiff := alt1 - alt2
	return math.Sqrt(dist*dist + altDiff*altDiff)
}

func haversine(latitude1, longitude1, latitude2, longitude2 float64) float64 {
	r := 6371009.0
	deg2rad := math.Pi / 180.0
	lat1 := latitude1 * deg2rad
	lat2 := latitude2 * deg2rad
	sin_dlat := math.Sin(math.Abs(lat1-lat2) * 0.5)
	sin_dlon := math.Sin(math.Abs(longitude1-longitude2) * deg2rad * 0.5)
	dist := r * 2.0 * math.Asin(math.Sqrt(
		sin_dlat*sin_dlat+math.Cos(lat1)*math.Cos(lat2)*sin_dlon*sin_dlon,
	))
	return dist
}
