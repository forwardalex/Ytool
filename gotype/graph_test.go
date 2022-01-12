package gotype

import (
	"fmt"
	"testing"
)

//先开启一个grpc服务测试
func TestBFS(t *testing.T) {
	b := NewGraph()
	var meta [][]int
	meta = [][]int{
		{4, 3, 3},
		{3, 2, 5},
		{2, 5, 3},
		{5, 4, 1},
		{3, 5, 3},
	}
	b.FillGraph(meta)
	fmt.Println(b.Nodes)
	node := b.Nodes.Get(2)
	b.BFSGraph(node)
	b.DFSGraph(node)
}
