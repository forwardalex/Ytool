package gotype

import (
	"fmt"
	"reflect"
	"sync"
)

type Graph struct {
	Nodes *NodesMap
	Edges *Set
	sync.RWMutex
}

// NodesMap 为查找方便使用map替代数组
type NodesMap struct {
	M map[int]*GNode
	sync.RWMutex
}

func (m *NodesMap) Put(key int, gNode *GNode) {
	m.RLock()
	defer m.RUnlock()
	m.M[key] = gNode
}
func (m *NodesMap) Add(gNode *GNode) {
	m.RLock()
	defer m.RUnlock()
	m.M[gNode.Value] = gNode
}

func (m *NodesMap) Get(key int) *GNode {
	m.RLock()
	defer m.RUnlock()
	return m.M[key]
}

// Contains node是否存在
func (m *NodesMap) Contains(key int) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.M[key]
	if ok {
		return true
	} else {
		return false
	}
}

func NewNodesMap() *NodesMap {
	return &NodesMap{
		M: make(map[int]*GNode, 0),
	}
}

// GNode 图的node
type GNode struct {
	Value int //node的值
	In    int //出度
	Out   int //入度
	Nexts *NodesMap
	Edges []*GEdges
}

// GEdges 图的边
type GEdges struct {
	weight int //权重
	from   *GNode
	to     *GNode
}

// NewGNode 新node
func NewGNode(value int) *GNode {
	return &GNode{
		Value: value,
		Nexts: NewNodesMap(),
		Edges: make([]*GEdges, 0),
	}
}

func NewGEdges(weight int, from, to *GNode) *GEdges {
	return &GEdges{
		weight: weight,
		from:   from,
		to:     to,
	}
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: NewNodesMap(),
		Edges: NewSet(),
	}
}

// FillGraph 将矩阵类型图边输入图中
/*
example
[][]int{  {1,2,3},   //1为from node  2为to node  3为weight
		  {2,5,7}
}
*/
func (g *Graph) FillGraph(meta [][]int) {
	for i := 0; i < len(meta); i++ {
		from := meta[i][0]
		to := meta[i][1]
		weight := meta[i][2]
		if !g.Nodes.Contains(from) {
			g.Nodes.Put(from, NewGNode(from))
		}
		if !g.Nodes.Contains(to) {
			g.Nodes.Put(to, NewGNode(to))
		}
		fromNode := g.Nodes.Get(from)
		toNode := g.Nodes.Get(to)
		newEdge := NewGEdges(weight, fromNode, toNode)
		fromNode.Nexts.Add(toNode)
		fromNode.Out++
		toNode.In++
		fromNode.Edges = append(fromNode.Edges, newEdge)
		g.Edges.Add(newEdge)
	}
}

//todo  dfn  bfs
func (g *Graph) BFSGraph(node *GNode) {
	if reflect.DeepEqual(node, &GNode{}) {
		return
	}
	queue := NewSliceQueue()
	hset := NewSet()
	queue.EnQueue(node)
	hset.Add(node)
	for !queue.IsEmpty() {
		cur := queue.DeQueue().(*GNode)
		fmt.Println("bfs ", cur.Value)
		for _, v := range cur.Nexts.M {
			if !hset.Contains(v) {
				hset.Add(v)
				queue.EnQueue(v)
			}
		}
	}
}

func (g *Graph) DFSGraph(node *GNode) {
	if reflect.DeepEqual(node, &GNode{}) {
		return
	}
	stack := NewSliceStack()
	hset := NewSet()
	stack.Push(node)
	hset.Add(node)
	fmt.Println("dfs ", node.Value)
	for !stack.IsEmpty() {
		cur := stack.Pop().(*GNode)
		for _, next := range cur.Nexts.M {
			if !hset.Contains(next) {
				stack.Push(cur)
				stack.Push(next)
				hset.Add(next)
				fmt.Println("dfs", next.Value)
				break
			}
		}
	}
}
