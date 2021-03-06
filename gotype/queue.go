package gotype

import "sync"

type Queue struct {
	arr []interface{}
	sync.RWMutex
}

func NewSliceQueue() *Queue {
	return &Queue{arr: make([]interface{}, 0)}
}

// Size 返回队列的大小
func (p *Queue) Size() int {
	return len(p.arr)
}

// IsEmpty 判断队列是否为空
func (p *Queue) IsEmpty() bool {
	return p.Size() == 0
}

// GetFront 返回队列首元素
func (p *Queue) GetFront() interface{} {
	if p.IsEmpty() {
		return nil
	}
	return p.arr[0]
}

// GetBack 返回队列尾元素
func (p *Queue) GetBack() interface{} {
	if p.IsEmpty() {
		return nil
	}
	return p.arr[p.Size()-1]
}

// EnQueue  把新元素加入队列尾
func (p *Queue) EnQueue(item interface{}) {
	p.Lock()
	defer p.Unlock()
	p.arr = append(p.arr, item)
}

// DeQueue   弹出头元素
func (p *Queue) DeQueue() interface{} {
	p.Lock()
	defer p.Unlock()
	if len(p.arr) != 0 {
		first := p.arr[0]
		p.arr = p.arr[1:]
		return first
	} else {
		return nil
	}
}

// Remove 简单实现一个Remove
func (p *Queue) Remove(item interface{}) {
	p.Lock()
	defer p.Unlock()
	for k, v := range p.arr {
		if v == item {
			p.arr = append(p.arr[:k], p.arr[k+1:]...)
		}
	}
}

// List  队列中所有元素
func (p *Queue) List() []interface{} {
	return p.arr
}

type LinkedQueue struct {
	head *LNode
	end  *LNode
	sync.RWMutex
}

func NewLinkedQueue() *LinkedQueue {
	return &LinkedQueue{}
}

// IsEmpty 判断队列是否为空,如果为空返回true，否则返回false
func (p *LinkedQueue) IsEmpty() bool {
	return p.head == nil
}

// Size 获取栈中元素的个数
func (p *LinkedQueue) Size() int {
	size := 0
	node := p.head
	for node != nil {
		node = node.Next
		size++
	}
	return size
}

// EnQueue 入队列：把元素e加到队列尾
func (p *LinkedQueue) EnQueue(e interface{}) {
	p.Lock()
	defer p.Unlock()
	node := &LNode{Data: e}
	if p.head == nil {
		p.head = node
		p.end = node
	} else {
		p.end.Next = node
		p.end = node
	}
}

// DeQueue 出队列，删除队列首元素
func (p *LinkedQueue) DeQueue() interface{} {
	p.Lock()
	defer p.Unlock()
	if p.head == nil {
		return nil
	}
	res := p.head
	p.head = p.head.Next
	if p.head == nil {
		p.end = nil
	}
	return res
}

// GetFront 取得队列首元素
func (p *LinkedQueue) GetFront() interface{} {
	if p.head == nil {
		return nil
	}
	return p.head.Data
}

// GetBack 取得队列尾元素
func (p *LinkedQueue) GetBack() interface{} {
	if p.end == nil {
		return nil
	}
	return p.end.Data
}
