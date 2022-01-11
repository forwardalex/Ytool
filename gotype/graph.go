package gotype

import "sync"

type Graph struct {
	Nodes *NodesMap
	Edges *Set
	sync.RWMutex
}

// NodesMap 为查找方便使用map替代数组
type NodesMap struct {
	m map[int]*GNode
	sync.RWMutex
}

func (m *NodesMap) Put(key int, gNode *GNode) {
	m.RLock()
	defer m.RUnlock()
	m.m[key] = gNode
}
func (m *NodesMap) Add(gNode *GNode) {
	m.RLock()
	defer m.RUnlock()
	m.m[gNode.value] = gNode
}

func (m *NodesMap) Get(key int) *GNode {
	m.RLock()
	defer m.RUnlock()
	return m.m[key]
}

// Contains node是否存在
func (m *NodesMap) Contains(key int) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.m[key]
	if ok {
		return true
	} else {
		return false
	}
}

func NewNodesMap() *NodesMap {
	return &NodesMap{
		m: make(map[int]*GNode, 0),
	}
}

// GNode 图的node
type GNode struct {
	value int //node的值
	in    int //出度
	out   int //入度
	nexts *NodesMap
	edges []*GEdges
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
		value: value,
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
[][]int{  1,2,3   //1为from node  2为to node  3为weight
		  2,5,7
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
		fromNode.nexts.Add(fromNode)
		fromNode.out++
		toNode.in++
		fromNode.edges = append(fromNode.edges, newEdge)
		g.Edges.Add(newEdge)
	}
}
