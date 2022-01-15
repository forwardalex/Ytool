package gotype

// ArrToSearchTree 有序数组处理成搜索二叉树
func ArrToSearchTree(arr []int, start, end int) *BNode {
	root := &BNode{}
	if start <= end {
		mid := (start + end + 1) >> 1
		root.Data = arr[mid]
		root.LeftChild = ArrToSearchTree(arr, start, mid-1)
		root.RightChild = ArrToSearchTree(arr, mid+1, end)
	} else {
		return nil
	}
	return root
}

// TreeByLevel 按层处理二叉树
func TreeByLevel(root *BNode, fn func(root *BNode)) {
	if root == nil {
		return
	}
	q := NewSliceQueue()
	q.EnQueue(root)
	for !q.IsEmpty() {
		root = q.DeQueue().(*BNode)
		fn(root)
		if root.LeftChild != nil {
			q.EnQueue(root.LeftChild)
		}
		if root.RightChild != nil {
			q.EnQueue(root.RightChild)
		}
	}
}
