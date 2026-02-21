package pool

import "sync"

// Resettable is an interface for types with Reset method.
type Resettable interface {
	Reset()
}

// Pool is a generic pool for Resettable types.
type Pool[T Resettable] struct {
	pool *sync.Pool
}

// New creates a new Pool instance.
func New[T Resettable](newFn func() T) *Pool[T] {
	return &Pool[T]{
		pool: &sync.Pool{
			New: func() any {
				return newFn()
			},
		},
	}
}

// Get retrieves an object from pool.
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put returns an object to pool after resetting it.
func (p *Pool[T]) Put(x T) {
	x.Reset()
	p.pool.Put(x)
}
