package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Durelius/next-week/internal/graph"
	"github.com/gorilla/mux"
)

func main() {
	slGraph, err := graph.NewWithData()
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/stopbyname/{name}", GetStopsByNameEndpoint).Methods("GET")
	r.HandleFunc("/path/{from}/{to}/{time}", GetPathEndpoint).Methods("GET")
	log.Println("Starting server at port 8080")
	http.ListenAndServe(":8080", corsMiddleware(r))
	log.Println("test")
	filteredVertices := slGraph.FindStopsByName("upplands väsby station")

	chosenStartPoint := filteredVertices[0]

	filteredVertices = slGraph.FindStopsByName("solna station")
	chosenDestination := filteredVertices[0]

	path := slGraph.FindRoute(chosenStartPoint, chosenDestination, 500)
	for _, edge := range path {
		log.Printf("TripID: %s, From: %s, To: %s, Start: %d, Arrival: %d", edge.Metadata.TripID, edge.SourceMethod(), edge.Destination(), edge.Metadata.Departure, edge.Metadata.Arrival)
	}

}
func GetStopsByNameEndpoint(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	nodes := graph.Instance().FindStopsByName(name)
	stops := []graph.Stop{}
	for _, v := range nodes {
		stops = append(stops, *v.Metadata())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(stops)
}
func GetPathEndpoint(w http.ResponseWriter, r *http.Request) {
	fromStopID := mux.Vars(r)["from"]
	toStopID := mux.Vars(r)["to"]
	startTimeStr := mux.Vars(r)["time"]
	startTimeHours, err := strconv.Atoi(startTimeStr[0:2])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	startTimeMinutes, err := strconv.Atoi(startTimeStr[3:5])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	minutesSinceMidnight := startTimeHours * 60
	minutesSinceMidnight += startTimeMinutes
	from := graph.Instance().GetVertexByID(fromStopID)
	to := graph.Instance().GetVertexByID(toStopID)
	path := graph.Instance().FindRoute(from, to, minutesSinceMidnight)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(path)
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
