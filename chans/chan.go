package chans

import (
	"golang.org/x/sync/errgroup"
)

type Result[T any] struct {
	Value T
	Err   error
}

type Chan[T any] chan Result[T]

func (c Chan[T]) Recv() (T, error) {
	r := <-c
	return r.Value, r.Err
}

func (c Chan[T]) Do(f func() (T, error)) {
	var r Result[T]
	r.Value, r.Err = f()
	c <- r
}

func (c Chan[T]) Go(g *errgroup.Group, f func() (T, error)) {
	g.Go(func() error {
		var r Result[T]
		r.Value, r.Err = f()
		c <- r
		return r.Err
	})
}

func Go[T any](g *errgroup.Group, f func() (T, error)) Chan[T] {
	c := make(Chan[T], 1)
	c.Go(g, f)
	return c
}
