package gotype

import (
	"fmt"
	"testing"
)

func TestArrToSearchTree(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	root := ArrToSearchTree(data, 0, len(data)-1)
	s := PrintSlice(root)
	fmt.Println("==", s)
}

//var s []interface{}
func PrintTreeByMid(root *BNode, s *[]interface{}) {
	if root == nil || root.Data == nil {
		return
	}
	PrintTreeByMid(root.LeftChild, s)
	*s = append(*s, root.Data)
	PrintTreeByMid(root.RightChild, s)
}

func PrintSlice(root *BNode) []interface{} {
	s := make([]interface{}, 0)
	PrintTreeByMid(root, &s)
	return s
}

func TestTreeByLevel(T *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	root := ArrToSearchTree(data, 0, len(data)-1)
	TreeByLevel(root, func(root *BNode) {
		fmt.Println("==", root.Data)
	})

}
