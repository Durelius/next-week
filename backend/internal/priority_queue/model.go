package priority_queue

// Package priority_queue defines a custom priorty queue for SL stops to be used in A* search.
// Implementation inspired by https://medium.com/@amankumarcs/priority-queue-in-go-b0b0b4844c91

type Item struct {
	value string // The actual data or value of the item
	g     int    // arrival time
	f     int    //  The priority of the item (lower value means higher priority)
	index int    // The index of the item in the heap (needed by the heap interface)
}

func NewItem(stopID string, g, f int) *Item {
	return &Item{
		value: stopID,
		f:     f,
		g:     g,
	}
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want the lowest priority (smallest integer) as the highest priority
	return pq[i].f < pq[j].f
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
func (i *Item) Value() string {
	return i.value
}

func (i *Item) G() int {
	return i.g
}

func (i *Item) F() int {
	return i.f
}
