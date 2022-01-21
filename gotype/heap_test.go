package gotype

import (
	"fmt"
	"testing"
)

func TestHeap(t *testing.T) {
	data := []int{6, 4, 2, 3, 10, 23, 1, 3, 6, 44, 2}
	maxh := NewMaxheap()
	for _, v := range data {
		maxh.HeapInsert(v)
	}
	fmt.Println(maxh.Arr)
	maxh.Heapifiy()
	fmt.Println("max heap", maxh.Arr)

	minh := NewMinheap()
	for _, v := range data {
		minh.HeapInsert(v)
	}
	fmt.Println(minh.Arr)
	minh.Heapifiy()
	fmt.Println("min heap", minh.Arr)

	data = []int{3, 1, 2, 9, 10, 11, 12}
	maxh = NewMaxheap()
	maxh.Arr = data
	maxh.Size = len(data)
	maxh.Heapifiy()

	fmt.Println("====", maxh.Arr)

}
