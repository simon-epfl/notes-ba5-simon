package heapStruct

type KthInt struct {
	Val uint64
	Idx int
}

// First time properly using AlgoDat :)
// use minheap for k-way merge like linux command (kinda like mergesort)
type MinHeap []KthInt

func (h MinHeap) Less(i int, j int) bool {
	a := h[i].Val < h[j].Val
	return a

}

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Swap(i int, j int) {
	a := h[i]
	h[i] = h[j]
	h[j] = a
}

func (h *MinHeap) Push(val interface{}) {
	*h = append(*h, val.(KthInt))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	size := len(old)
	val := old[size-1]
	*h = old[:size-1]
	return val
}
