package atomicid

import "sync/atomic"

type AtomicId uint64

func New(id *uint64) *AtomicId {
	return (*AtomicId)(id)
}

func (p *AtomicId) Next() uint64 {
	return atomic.AddUint64((*uint64)(p), 1)
}

func (p *AtomicId) Current() uint64 {
	return atomic.LoadUint64((*uint64)(p))
}
