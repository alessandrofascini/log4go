package tests

import "math"

type PriorityQueueComparator func(a int, b int) bool

type PriorityQueue struct {
	heap     []int
	capacity int
	length   int
	compare  PriorityQueueComparator
	height   int
}

func NewPriorityQueue(capacity uint, compare PriorityQueueComparator) *PriorityQueue {
	height := int(math.Log2(float64(capacity))) + 1
	return &PriorityQueue{
		heap:     make([]int, 1<<height),
		capacity: height,
		length:   0, // first available cell
		compare:  compare,
		height:   height,
	}
}

// Push add an element to a priority queue
func (pq *PriorityQueue) Push(e int) {
	if pq.length == pq.capacity {
		pq.heap = append(pq.heap, make([]int, pq.height)...)
		pq.height = pq.height << 1
	}
	pq.heap[pq.length] = e
	pq.moveOn(pq.length)
	pq.length++
}

// Peek get the first element, but not pop it
func (pq *PriorityQueue) Peek() int {
	// TODO Better here
	return pq.heap[0]
}

// Pop get the first element of the queue and then drop it
func (pq *PriorityQueue) Pop() int {
	firstElement := pq.heap[0]
	pq.moveDown(0)
	return firstElement
}

func (pq *PriorityQueue) moveOn(i int) {
	parent := (i - 1) / 2
	for i != 0 && pq.compare(pq.heap[i], pq.heap[parent]) {
		pq.heap[i], pq.heap[parent] = pq.heap[parent], pq.heap[i]
		i = parent
		parent = (i - 1) / 2
	}
}

// TODO Review this
func (pq *PriorityQueue) moveDown(i int) {
	child1 := 1<<i + 1
	child2 := 1<<i + 2
	for child2 <= pq.length {
		if pq.compare(child1, child2) {
			pq.heap[i] = pq.heap[child1]
			i = child1
		} else {
			pq.heap[i] = pq.heap[child2]
			i = child2
		}
		child1 = 1<<i + 1
		child2 = 1<<i + 2
	}
	if child1 <= pq.length {
		pq.heap[i] = pq.heap[child1]
	} else {
		pq.length--
	}
}
