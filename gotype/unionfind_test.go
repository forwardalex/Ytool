package gotype

import (
	"fmt"
	"testing"
)

func TestUnionSet(t *testing.T) {
	set := NewUnionFind(10)
	set.Union(1, 2)
	set.Union(1, 3)
	set.Union(2, 3)
	set.Union(2, 4)
	set.Union(2, 3)
	is := set.IsConnected(1, 2)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(set.set)
	fmt.Println(is)
	//fmt.Println(set.Root(10))
	//fmt.Println(set.size[-1])
}
