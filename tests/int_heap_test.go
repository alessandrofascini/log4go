package tests

import (
	"container/heap"
	"fmt"
	"testing"
)

// This example inserts several ints into an IntHeap, checks the minimum,
// and removes them in order of priority.
func TestIntHeap(t *testing.T) {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minimum: %d\n", (*h)[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
}
