package gotype

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	s := NewSliceStack()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	for !s.IsEmpty() {
		fmt.Println(s.Pop())
	}
}
