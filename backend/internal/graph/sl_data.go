package graph

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/gocarina/gocsv"
)

func (graph *SLGraph) init() error {

	_, _, stopTimes, stops, _, err := load()
	if err != nil {
		return err
	}
	for _, stop := range stops {
		v := NewVertex(stop.StopID)
		stop.StopNameLower = strings.ToLower(stop.StopName)
		v.SetMetadata(stop)
		graph.AddVertex(v)
	}
	stopTimeMap := make(map[string][]StopTimes)
	for _, stopTime := range stopTimes {
		stopTimeMap[stopTime.TripID] = append(stopTimeMap[stopTime.TripID], *stopTime)
	}
	for tripID, times := range stopTimeMap {

		sort.Slice(times, func(i, j int) bool {
			return times[i].StopSequence < times[j].StopSequence
		})

		stopTimeMap[tripID] = times
	}
	for _, times := range stopTimeMap {
		for i := 0; i < len(times)-1; i++ {
			from := times[i]
			to := times[i+1]
			fromVertice := graph.GetVertexByID(from.StopID)
			toVertice := graph.GetVertexByID(to.StopID)
			edgeProps := EdgeProperties{
				TripID:         from.TripID,
				Departure:      toMinutes(from.DepartureTime),
				Arrival:        toMinutes(to.ArrivalTime),
				TransferType:   COMMUTE_EDGE,
				SourceStopName: fromVertice.metadata.StopName,
				DestStopName:   toVertice.metadata.StopName,
			}
			if _, err := graph.AddEdge(fromVertice, toVertice, edgeProps); err != nil {
				return err
			}
		}
	}
	return graph.addTransferEdges(stops)
}
func (graph *SLGraph) addTransferEdges(stops []*Stop) error {
	for i, a := range stops {
		for j, b := range stops {
			if i >= j {
				continue
			} // skip self and avoid duplicates
			if a.StopName == b.StopName {
				continue
			}
			dist, err := ApproxDistanceMeters(a, b)
			if err != nil {
				log.Printf("Couldn't calculate distance, err: %v", err)
				continue
			}

			if dist < 400 {
				// 80 meters per minute
				walkMinutes := int(dist / 80.0)
				//minimum of 1 minute
				if walkMinutes == 0 {
					walkMinutes = 1
				}
				edgeProps := EdgeProperties{
					Departure:    0,
					Arrival:      walkMinutes,
					TransferType: WALK_EDGE,
				}
				from := graph.GetVertexByID(a.StopID)
				to := graph.GetVertexByID(b.StopID)
				edgeProps.SourceStopName = from.metadata.StopName
				edgeProps.DestStopName = to.metadata.StopName

				if _, err := graph.AddEdge(from, to, edgeProps); err != nil {
					return err
				}
				if _, err := graph.AddEdge(to, from, edgeProps); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func load() ([]*Agency, []*Routes, []*StopTimes, []*Stop, []*Trips, error) {
	var agencies []*Agency
	file, err := os.Open(PATH_AGENCY)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer file.Close()
	if err := gocsv.Unmarshal(file, &agencies); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var routes []*Routes
	file, err = os.Open(PATH_ROUTES)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer file.Close()
	if err := gocsv.Unmarshal(file, &routes); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var stopTimes []*StopTimes
	file, err = os.Open(PATH_STOPTIMES)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer file.Close()
	if err := gocsv.Unmarshal(file, &stopTimes); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var stops []*Stop
	file, err = os.Open(PATH_STOPS)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer file.Close()
	if err := gocsv.Unmarshal(file, &stops); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var trips []*Trips
	file, err = os.Open(PATH_TRIPS)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer file.Close()
	if err := gocsv.Unmarshal(file, &trips); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return agencies, routes, stopTimes, stops, trips, err
}
func toMinutes(t string) int {
	var h, m, s int
	fmt.Sscanf(t, "%d:%d:%d", &h, &m, &s)
	return h*60 + m
}
