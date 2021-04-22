package clist

import "sync/atomic"

const marked = 1

type bitflag struct {
	data uint32
}

// setMarked set the node to marked status.
func (b *bitflag) setMarked() {
	for {
		old := atomic.LoadUint32(&b.data)
		if old&marked != marked {
			// make sure do CAS atomic successfully
			if atomic.CompareAndSwapUint32(&b.data, old, marked) {
				return
			}
			continue
		}
		return
	}
}

// isMarked check the node whether is marked.
func (b *bitflag) isMarked() bool {
	return atomic.LoadUint32(&b.data) == marked
}
