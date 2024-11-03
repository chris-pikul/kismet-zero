// Copied from github.com/kelindar/search, original copyright is as follows:
//
// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
package llama

import "io"

// Pool is a generic pool of resources that can be reused.
type pool[T io.Closer] struct {
	pool chan T
	make func() T
}

// newPool creates a new pool of resources.
func newPool[T io.Closer](size int, new func() T) *pool[T] {
	return &pool[T]{
		pool: make(chan T, size),
		make: new,
	}
}

// Get returns a resource from the pool or creates a new one.
func (p *pool[T]) Get() T {
	select {
	case x := <-p.pool:
		return x
	default:
		return p.make()
	}
}

// Put returns the resource to the pool.
func (p *pool[T]) Put(x T) {
	select {
	case p.pool <- x:
	default:
		x.Close() // Close the resource if the pool is full
	}
}

// Close closes the pool and releases any resources associated with it.
func (p *pool[T]) Close() {
	close(p.pool)
	for x := range p.pool {
		x.Close()
	}
}
