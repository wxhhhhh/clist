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

func NewInt() *IntList {
	return &IntList{head: newIntNode(0)}
}

func (l *IntList) getHead() *intNode {
	return l.head
}

func (l *IntList) findNode(value int) (*intNode, *intNode) {

	a := l.head
	b := a.loadNext()

	for b != nil && b.value < value {
		a = b
		b = b.loadNext()
	}

	return a, b
}

func (l *IntList) Insert(value int) bool {

	for {

		a, b := l.findNode(value)

		if b != nil && b.value == value {
			return false
		}

		a.mu.Lock()
		if a.loadNext() != b {
			a.mu.Unlock()
			continue
		}

		x := newIntNode(value)
		x.next = b

		a.storeNext(x)
		a.mu.Unlock()

		atomic.AddInt64(&l.length, 1)

		return true
	}

}

func (l *IntList) Delete(value int) bool {

	for {

		a, b := l.findNode(value)

		if b == nil || b.value != value {
			return false
		}

		b.mu.Lock()
		if b.flag.IsMarked() {
			b.mu.Unlock()
			continue
		}

		a.mu.Lock()
		if a.loadNext() != b || a.flag.IsMarked() {
			a.mu.Unlock()
			b.mu.Unlock()
			continue
		}

		b.flag.SetMarked()
		a.storeNext(b.next)
		a.mu.Unlock()
		b.mu.Unlock()
		atomic.AddInt64(&l.length, -1)
		return true
	}

}

func (l *IntList) Contains(value int) bool {

	x := l.head.loadNext()
	for x != nil && x.value < value {
		x = x.loadNext()
	}
	if x == nil {
		return false
	}
	return !x.flag.IsMarked() && x.value == value

}

func (l *IntList) Range(f func(value int) bool) {
	x := l.head.loadNext()
	for x != nil {
		if x.flag.IsMarked() {
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
