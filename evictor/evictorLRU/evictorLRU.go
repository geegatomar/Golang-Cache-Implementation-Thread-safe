package evictorLRU

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	ErrorNotFound = errors.New("Error: Element with this key not found in cache.")
)

// EvictorLRU implements the Evictor interface
type EvictorLRU struct {
	pq PriorityQueue
}

func (e *EvictorLRU) Evict(evictElement func(key string) bool) {
	for {
		// Typecasting the empty interface type to struct
		if len(e.pq) == 0 {
			log.Panicf("Length of pq has become 0")
			return
		}
		popped := heap.Pop(&e.pq).(*element)
		fmt.Printf("Popping ----------- %v", popped)
		ele := popped.key
		if !evictElement(ele) {
			return
		}
	}
}

func (e *EvictorLRU) OnAdd(key string) error {
	ele := element{
		key:              key,
		timeOfLastAccess: time.Now(),
		priority:         time.Now().UnixNano(),
	}
	heap.Push(&e.pq, &ele)
	return nil
}

func (e *EvictorLRU) OnAccess(key string) error {
	for _, ele := range e.pq {
		if ele.key == key {
			e.pq.update(ele, time.Now(), time.Now().UnixNano())
			return nil
		}
	}
	return ErrorNotFound
}

func (e *EvictorLRU) OnDelete(key string) error {
	// Idea here is to pop until you find the right one, storing them temporarily in temp,
	// and then push those popped ones back in.
	temp := make([]*element, 0)

	// Pop elements until you find the one
	for e.pq.Len() > 0 {
		popped := heap.Pop(&e.pq).(*element)
		ele := popped.key
		if ele == key {
			break
		} else {
			temp = append(temp, popped)
		}
	}

	// Push back all the remaining elements
	for _, ele := range temp {
		heap.Push(&e.pq, ele)
	}
	return nil
}

func (e *EvictorLRU) OnUpdate(key string, value string) error {
	// Idea here is to pop until you find the right one, storing them temporarily in temp,
	// and then push those popped ones back in.
	temp := make([]*element, 0)

	// Pop elements until you find the one
	for e.pq.Len() > 0 {
		popped := heap.Pop(&e.pq).(*element)
		ele := popped.key
		if ele == key {
			// Main part here
			e.pq.update(popped, time.Now(), time.Now().UnixNano())

		} else {
			temp = append(temp, popped)
		}
	}

	// Push back all the remaining elements
	for _, ele := range temp {
		heap.Push(&e.pq, ele)
	}
	return nil
}
