package chans

import (
	"context"

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

func Go[T any](g *errgroup.Group, f func() (T, error)) Chan[T] {
	ch := make(Chan[T], 1)
	if g == nil {
		go func() {
			var r Result[T]
			r.Value, r.Err = f()
			ch <- r
		}()
	} else {
		g.Go(func() error {
			var r Result[T]
			r.Value, r.Err = f()
			ch <- r
			return r.Err
		})
	}
	return ch
}

func Recv[T any](ctx context.Context, ch Chan[T]) (T, error) {
	select {
	case r := <-ch:
		return r.Value, r.Err
	case <-ctx.Done():
		var z T
		return z, ctx.Err()
	}
}
