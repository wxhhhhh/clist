package clist

import (
	"sync"
	"sync/atomic"
)

type IntList struct {
	head   *intNode
	length int64
	mu     sync.RWMutex
}

// IntList return a dummy head of clist.
func NewInt() *IntList {
	return &IntList{head: newIntNode(0)}
}

// findNodes return node a and b which satisfy a.value < value <= b.value.
func (l *IntList) findNodes(value int) (*intNode, *intNode) {

	a := l.head
	b := a.loadNext()

	for b != nil && b.value < value {
		a = b
		b = b.loadNext()
	}

	return a, b
}

// Insert add value in to the clist, return true if insert node successfully.
func (l *IntList) Insert(value int) bool {

	for {
		// 1. find node a and b, if value is already in clist return false
		a, b := l.findNodes(value)
		if b != nil && b.value == value {
			return false
		}

		// 2. lock node a, check whether a.next != b
		// if so, unlock node a, continue step1
		a.mu.Lock()
		if a.loadNext() != b {
			a.mu.Unlock()
			continue
		}

		// 3. create new node x
		x := newIntNode(value)

		// 4. x.next = b, a.next = x
		x.next = b
		a.storeNext(x)

		// 5. unlock node a
		a.mu.Unlock()

		// 6. the length of clist plus one
		atomic.AddInt64(&l.length, 1)

		return true
	}

}

// Delete delete the node from clist, return true if delete successfully.
func (l *IntList) Delete(value int) bool {

	for {
		// 1. find node a and b, if value is not in clist return false
		a, b := l.findNodes(value)
		if b == nil || b.value != value {
			return false
		}

		// 2. lock node b, check whether b is marked
		// if so, unlock b, continue step1
		b.mu.Lock()
		if b.flag.isMarked() {
			b.mu.Unlock()
			continue
		}

		// 3. lock node a, check whether a.next !=b or a is marked
		// if so, unblock a and b, continue step1
		a.mu.Lock()
		if a.loadNext() != b || a.flag.isMarked() {
			a.mu.Unlock()
			b.mu.Unlock()
			continue
		}

		// 4. mark b, set a.next = b.next
		b.flag.setMarked()
		a.storeNext(b.next)

		// 5. unlock node a and node b
		a.mu.Unlock()
		b.mu.Unlock()

		// 6. the length of clist minus one
		atomic.AddInt64(&l.length, -1)
		return true
	}

}

// Contains check if value is in clist.
func (l *IntList) Contains(value int) bool {

	x := l.head.loadNext()
	for x != nil && x.value < value {
		x = x.loadNext()
	}
	// find x, and check whether the node is still in clist or not.
	if x != nil && x.value == value {
		return !x.flag.isMarked()
	}
	return false

}

// Range retrieve all elements in clist and call func f for each.
func (l *IntList) Range(f func(value int) bool) {

	x := l.head.loadNext()
	for x != nil {
		// ignore the marked node.
		if x.flag.isMarked() {
			x = x.loadNext()
			continue
		}
		if !f(x.value) {
			break
		}
		x = x.loadNext()
	}
}

// Len return the length of list atomically.
func (l *IntList) Len() int {
	return int(atomic.LoadInt64(&l.length))
}
