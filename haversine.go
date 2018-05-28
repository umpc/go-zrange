package zrange

import "math"

// Haversine returns the great circle distance between two latitude/longitude
// points in kilometers.
//
// Derived from: https://web.archive.org/web/20180528040024/https://community.esri.com/groups/coordinate-reference-systems/blog/2017/10/05/haversine-formula
//
func Haversine(lat1, lng1, lat2, lng2 float64) float64 {
	phi1 := degToRad(lat1)
	phi2 := degToRad(lat2)

	deltaPhi := degToRad(lat2 - lat1)
	deltaLambda := degToRad(lng2 - lng1)

	a := hav(deltaPhi) + math.Cos(phi1)*math.Cos(phi2)*hav(deltaLambda)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * earthSemiMajorAxis
}

func degToRad(x float64) float64 {
	const degToRad = math.Pi / 180
	return x * degToRad
}

func hav(x float64) float64 {
	return math.Pow(math.Sin(x/2), 2)
}
