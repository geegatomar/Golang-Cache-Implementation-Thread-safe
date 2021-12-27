package evictorLRU

import (
	"container/heap"
	"time"
)

// Check: https://pkg.go.dev/container/heap#example-package-PriorityQueue

// Creating a MinHeap where the element which has lowest priority (least recently used element) will be popped first
type element struct {
	key              string
	timeOfLastAccess time.Time
	priority         int64
	index            int
}

// A PriorityQueue implements heap.Interface and holds pointers to 'element'.
type PriorityQueue []*element

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest time, i.e. the one used the least recently
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	ele := x.(*element)
	ele.index = n
	*pq = append(*pq, ele)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	ele := old[n-1]
	old[n-1] = nil // avoid memory leak
	ele.index = -1 // for safety
	*pq = old[0 : n-1]
	return ele
}

func (pq *PriorityQueue) update(ele *element, latest time.Time, priority int64) {
	ele.timeOfLastAccess = latest
	ele.priority = priority
	heap.Fix(pq, ele.index)
}
