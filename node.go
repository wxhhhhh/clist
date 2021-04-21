package clist

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type intNode struct {
	value int
	mu    sync.Mutex
	next  *intNode
	flag  bitflag
}

// newIntNode create a node with value v.
func newIntNode(v int) *intNode {
	return &intNode{
		value: v,
	}
}

// loadNext return n.next atomically.
func (n *intNode) loadNext() *intNode {
	return (*intNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.next))))
}

// storeNext set n.next = node atomically.
func (n *intNode) storeNext(node *intNode) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&n.next)), unsafe.Pointer(node))
}
