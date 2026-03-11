package graph

import (
	"container/heap"
	"log"
	"slices"
	"sort"
	"strings"

	pq "github.com/Durelius/next-week/internal/priority_queue"
)

// FindRoute finds the fastest way between two stops using a custom implementation of the A* algorithm
func (graph *SLGraph) FindRoute(start *Vertex, destination *Vertex, startTime int) []*Edge {
	open := make(pq.PriorityQueue, 0)
	heap.Init(&open)
	heap.Push(&open, pq.NewItem(start.metadata.StopID, startTime, startTime))
	closed := make(map[string]bool)
	bestG := make(map[string]int)
	bestG[start.label] = startTime
	cameFrom := make(map[string]*Edge) //"To reach this stop, we used this edge"
	for len(open) > 0 {
		current := heap.Pop(&open).(*pq.Item)
		currentStop := graph.GetVertexByID(current.Value())

		currentTripID := ""
		if prevEdge, ok := cameFrom[current.Value()]; ok {
			currentTripID = prevEdge.Metadata.TripID
		}
		if current.Value() == destination.label {

			var path []*Edge
			currentID := destination.label

			//find starting point
			for currentID != start.label {
				edge, ok := cameFrom[currentID]
				if !ok {
					log.Println("error: Didn't find came from")
					return nil
				}

				path = append(path, edge)
				currentID = edge.source.label // go backwards
			}

			slices.Reverse(path)
			return path
		}

		if closed[current.Value()] {
			continue
		}

		closed[current.Value()] = true
		for _, edge := range currentStop.edges {
			neighborID := edge.dest.label
			newG := edge.calculateG(current.G(), currentTripID)
			if newG == -1 {
				continue
			}
			if best, exists := bestG[neighborID]; !exists || newG < best {
				bestG[neighborID] = newG
				neighborStop := graph.GetVertexByID(neighborID)
				h := calculateH(neighborStop.metadata, destination.metadata)
				f := newG + h
				heap.Push(&open, pq.NewItem(neighborID, newG, f))
				cameFrom[neighborID] = edge
			}
		}

	}
	return nil
}

func (graph *SLGraph) FindStopsByName(name string) []*Vertex {
	filteredVertices := []*Vertex{}
	name = strings.ToLower(name)
	for _, v := range graph.GetAllVertices() {
		if strings.Contains(v.metadata.StopNameLower, name) {
			filteredVertices = append(filteredVertices, v)
		}
	}
	sort.Slice(filteredVertices, func(i, j int) bool {
		return filteredVertices[i].label < filteredVertices[j].label
	})
	return filteredVertices
}
