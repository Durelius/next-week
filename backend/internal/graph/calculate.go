package graph

import (
	"math"
	"strconv"
)

func calculateH(from *Stop, destination *Stop) int {
	dist, err := ApproxDistanceMeters(from, destination)
	if err != nil {
		return 0
	}
	// tunnelbana top speed 70 km/h = 1166 meters per minute
	// this is heuristics, so not accurate just a guess
	return int(dist / 1166.0)
}
func (e *Edge) calculateG(currentTime int, currentTripID string) int {
	if e.Metadata.TransferType == WALK_EDGE {
		//penalty 5 here for walking
		return currentTime + e.Metadata.Arrival + 5
	}

	if e.Metadata.Departure < currentTime {
		return -1
	}
	waitTime := e.Metadata.Departure - currentTime
	travelTime := e.Metadata.Arrival - e.Metadata.Departure
	penalty := 0
	if currentTripID != "" && e.Metadata.TripID != "" && currentTripID != e.Metadata.TripID {
		//penalty 5 here for changing line
		penalty = 5
	}
	return currentTime + waitTime + travelTime + penalty
}

// ApproxDistanceMeters calculates the distance in meters between to stops  by their coordinates
// algorithm implementation found at https://github.com/daveroberts0321/distancecalculator/blob/main/distancecalculator.go
func ApproxDistanceMeters(from *Stop, to *Stop) (float64, error) {
	fromLat, err := strconv.ParseFloat(from.StopLatitude, 64)
	if err != nil {
		return 0, err
	}
	fromLong, err := strconv.ParseFloat(from.StopLongitude, 64)
	if err != nil {
		return 0, err
	}
	toLat, err := strconv.ParseFloat(to.StopLatitude, 64)
	if err != nil {
		return 0, err
	}
	toLong, err := strconv.ParseFloat(to.StopLongitude, 64)
	if err != nil {
		return 0, err
	}
	// Calculate distances
	lat1 := fromLat * math.Pi / 180
	long1 := fromLong * math.Pi / 180
	r := 6378100.0 // Earth radius in METERS

	lat2 := toLat * math.Pi / 180
	long2 := toLong * math.Pi / 180

	// Haversine formula to calculate distance between two points
	h := math.Pow(math.Sin((lat2-lat1)/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin((long2-long1)/2), 2)

	return 2 * r * math.Asin(math.Sqrt(h)), nil
}
