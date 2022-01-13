package gotype

import (
	"fmt"
	"testing"
)

func TestHeap(t *testing.T) {
	data := []int{6, 4, 2, 3, 10, 99}
	maxh := NewMaxheap()
	for _, v := range data {
		maxh.HeapInsert(v)
	}
	fmt.Println(maxh.Arr)
	maxh.Heapifiy()
	fmt.Println(maxh.Arr)

	minh := NewMinheap()
	for _, v := range data {
		minh.HeapInsert(v)
	}
	fmt.Println(minh.Arr)
	minh.Heapifiy()
	fmt.Println(minh.Arr)
}
