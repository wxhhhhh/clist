package clist

import (
	"sync"
)

type bitflag struct {
	mu     sync.RWMutex
	marked bool
}

// setMarked set the node to marked status.
func (b *bitflag) setMarked() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.marked = true
}

// isMarked check the node whether is marked.
func (b *bitflag) isMarked() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.marked
}
