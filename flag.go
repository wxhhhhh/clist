package clist

import "sync/atomic"

const marked = 1

type bitflag struct {
	data uint32
}

func (b *bitflag) SetMarked() {
	for {
		old := atomic.LoadUint32(&b.data)
		if old&marked != marked {
			if atomic.CompareAndSwapUint32(&b.data, old, marked) {
				return
			}
			continue
		}
		return
	}
}

func (b *bitflag) IsMarked() bool {
	return atomic.LoadUint32(&b.data) == marked
}
