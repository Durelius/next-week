package graph

import (
	"log"
	"sync"
	"sync/atomic"
)

type SLGraph struct {
	vertices      map[string]*Vertex
	edges         map[string]map[string]*Edge // source -> dest -> edge
	verticesCount uint32
	edgesCount    uint32
}

// New creates a New empty SL graph
func New() *SLGraph {
	return &SLGraph{
		vertices: make(map[string]*Vertex),
		edges:    make(map[string]map[string]*Edge),
	}
}

var (
	instance *SLGraph
	once     sync.Once
)

// NewWithData creates a non-empty SL graph filled with data from hard-coded .CSV files
func NewWithData() (*SLGraph, error) {
	var initErr error
	once.Do(func() {
		graph := New()
		log.Println("loading graph data....")
		if err := graph.init(); err != nil {
			initErr = err
			return
		}
		instance = graph
	})
	if initErr != nil {
		return nil, initErr
	}
	return instance, nil
}
func Instance() *SLGraph {
	if instance == nil {
		log.Fatal("graph is nil, call NewWithData first")
	}
	return instance
}
func (graph *SLGraph) Order() uint32 {
	return atomic.LoadUint32(&graph.verticesCount)
}

func (graph *SLGraph) Size() uint32 {
	return atomic.LoadUint32(&graph.edgesCount)
}
