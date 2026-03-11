package graph

import (
	"sync/atomic"
)

// ----------------------------
// Edge struct
// ----------------------------

type EdgeType int

// EDGE TYPES
const (
	WALK_EDGE EdgeType = iota + 1
	COMMUTE_EDGE
)

type EdgeProperties struct {
	TripID         string   `json:"tripId"`
	Departure      int      `json:"departure"` // minutes since midnight
	Arrival        int      `json:"arrival"`
	TransferType   EdgeType `json:"transferType"`
	SourceStopName string   `json:"sourceStopName"`
	DestStopName   string   `json:"destStopName"`
}

type Edge struct {
	source   *Vertex
	dest     *Vertex
	Metadata EdgeProperties `json:"metadata"`
}

func NewEdge(from, to *Vertex, metadata EdgeProperties) *Edge {
	return &Edge{source: from, dest: to, Metadata: metadata}
}
func (graph *SLGraph) AddEdge(from, to *Vertex, metadata EdgeProperties) (*Edge, error) {
	if from == nil || to == nil {
		return nil, ErrNilVertices
	}

	// add vertices if they don't exist
	if _, ok := graph.vertices[from.label]; !ok {
		graph.AddVertex(from)
	}
	if _, ok := graph.vertices[to.label]; !ok {
		graph.AddVertex(to)
	}

	edge := NewEdge(from, to, metadata)

	// append to outgoing edges
	from.AddEdge(edge)

	// add to edge map
	if _, ok := graph.edges[from.label]; !ok {
		graph.edges[from.label] = map[string]*Edge{to.label: edge}
	} else {
		graph.edges[from.label][to.label] = edge
	}

	atomic.AddUint32(&graph.edgesCount, 1)
	return edge, nil
}
func (e *Edge) Source() *Stop {
	return e.source.metadata
}
func (e *Edge) Destination() *Stop {
	return e.dest.metadata
}

func (graph *SLGraph) EdgesOf(v *Vertex) []*Edge {
	if v == nil || graph.vertices[v.label] == nil {
		return nil
	}
	edges := make([]*Edge, 0)
	// outgoing edges
	if destMap, ok := graph.edges[v.label]; ok {
		for _, edge := range destMap {
			edges = append(edges, edge)
		}
	}
	// incoming edges
	for srcLabel, destMap := range graph.edges {
		if srcLabel == v.label {
			continue
		}
		if edge, ok := destMap[v.label]; ok {
			edges = append(edges, edge)
		}
	}
	return edges
}
func (graph *SLGraph) AllEdges() []*Edge {
	var all []*Edge
	for _, destMap := range graph.edges {
		for _, edge := range destMap {
			all = append(all, edge)
		}
	}
	return all
}
func (graph *SLGraph) ContainsEdge(from, to *Vertex) bool {
	if from == nil || to == nil {
		return false
	}
	if destMap, ok := graph.edges[from.label]; ok {
		if _, ok2 := destMap[to.label]; ok2 {
			return true
		}
	}
	return false
}

// ----------------------------
// Remove methods
// ----------------------------

func (graph *SLGraph) RemoveEdges(edges ...*Edge) {
	for _, e := range edges {
		if e == nil || graph.vertices[e.source.label] == nil || graph.vertices[e.dest.label] == nil {
			continue
		}
		delete(graph.edges[e.source.label], e.dest.label)
		// remove from outgoing edges slice
		src := graph.vertices[e.source.label]
		for i, edge := range src.edges {
			if edge.dest.label == e.dest.label {
				src.edges = append(src.edges[:i], src.edges[i+1:]...)
				break
			}
		}
		atomic.AddUint32(&graph.edgesCount, ^(uint32(1) - 1))
		e.dest.inDegree--
	}
}

func removeEdgeFromSlice(edges []*Edge, target *Edge) []*Edge {
	for i, e := range edges {
		if e == target {
			return append(edges[:i], edges[i+1:]...)
		}
	}
	return edges
}
