package service

import "math"

const (
	earthRadiusKm = 6371.0
	baseFare      = 10000.0 // VND
	pricePerKm    = 5000.0  // VND
)

// EstimateFare returns distance in km and estimated fare in VND.
func EstimateFare(pickupLat, pickupLng, dropoffLat, dropoffLng float64) (distanceKm float64, fare float64) {
	distanceKm = haversineKm(pickupLat, pickupLng, dropoffLat, dropoffLng)
	fare = baseFare + distanceKm*pricePerKm
	return distanceKm, math.Round(fare)
}

func haversineKm(lat1, lng1, lat2, lng2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}
