package gotype

type Minheap struct {
	Arr []int
	Size int
}
type Maxheap struct {
	Arr []int
	Size int
}

func NewMinheap() *Minheap{
	return &Minheap{}
}
func NewMaxheap ()*Maxheap{
	return &Maxheap{}
}

// HeapInsert 大跟堆生成
func (h *Maxheap)HeapInsert(value int)  {
	h.Size=h.Size+1
	h.Arr=append(h.Arr,value)
	pos:=h.Size-1
	root:=(pos-1)>>1
	for root>=0{
		if h.Arr[pos]>h.Arr[root]{
			h.Arr[pos],h.Arr[root]=h.Arr[root],h.Arr[pos]
			pos=root
			root=(root-1)>>2
		}else {
			break
		}
	}
}

// Heapifiy 调整成大根堆
func(h *Maxheap)Heapifiy(){
	root:=0
	Child:=2*root+1
	for Child<h.Size{
		if Child+1<h.Size && h.Arr[Child]<h.Arr[Child+1]{
			Child++
		}
		if h.Arr[Child]>h.Arr[root]{
			h.Arr[Child],h.Arr[root]=h.Arr[root],h.Arr[Child]
			root=Child
			Child=2*root+1
		}else {
			break
		}
	}
}

// HeapInsert 小根堆
func (h *Minheap)HeapInsert(value int){
	h.Size=h.Size+1
	h.Arr=append(h.Arr,value)
	pos:=h.Size-1
	root:=(pos-1)>>1
	for root>=0{
		if h.Arr[pos]<h.Arr[root]{
			h.Arr[pos],h.Arr[root]=h.Arr[root],h.Arr[pos]
			pos=root
			root=(root-1)>>2
		}else {
			break
		}
	}
}

// Heapifiy  小根堆调整
func  (h *Minheap)Heapifiy()  {
	root:=0
	Child:=2*root+1
	for Child<h.Size{
		if Child+1<h.Size && h.Arr[Child]>h.Arr[Child+1]{
			Child++
		}
		if h.Arr[Child]<h.Arr[root]{
			h.Arr[Child],h.Arr[root]=h.Arr[root],h.Arr[Child]
			root=Child
			Child=2*root+1
		}else {
			break
		}
	}
}

type Person struct {
	Age int
}


