package chans

import "context"

type Result[T any] struct {
	Value T
	Err   error
}

type Chan[T any] chan Result[T]

func (c Chan[T]) Recv() (T, error) {
	r := <-c
	return r.Value, r.Err
}

func Go[T any](ctx context.Context, f func() (T, error)) Chan[T] {
	ch := make(Chan[T])
	go func() {
		var r Result[T]
		r.Value, r.Err = f()
		select {
		case ch <- r:
		case <-ctx.Done():
		}
	}()
	return ch
}
