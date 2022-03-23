package chans

import (
	"context"

	"github.com/icholy/exp/slices"
)

func Merge[T any](ctx context.Context, chans ...Chan[T]) Chan[T] {
	ch := make(Chan[T])
	go merge(ctx, ch, chans)
	return ch
}

func merge[T any](ctx context.Context, out Chan[T], chans []Chan[T]) {
	switch len(chans) {
	case 1:
		merge1(ctx, out, chans)
	case 2:
		merge2(ctx, out, chans)
	case 3:
		merge3(ctx, out, chans)
	case 4:
		merge4(ctx, out, chans)
	case 5:
		merge5(ctx, out, chans)
	default:
		for _, batch := range slices.Batch(chans, 5) {
			go merge(ctx, out, batch)
		}
	}
}

func merge1[T any](ctx context.Context, out Chan[T], chans []Chan[T]) {
	for {
		var r Result[T]
		select {
		case r = <-chans[0]:
		case <-ctx.Done():
		return
		}
		select {
		case out <- r:
		case <-ctx.Done():
			return
		}
	}
}

func merge2[T any](ctx context.Context, out Chan[T], chans []Chan[T]) {
	for {
		var r Result[T]
		select {
		case r = <-chans[0]:
		case r = <-chans[1]:
		case <-ctx.Done():
		return
		}
		select {
		case out <- r:
		case <-ctx.Done():
			return
		}
	}
}

func merge3[T any](ctx context.Context, out Chan[T], chans []Chan[T]) {
	for {
		var r Result[T]
		select {
		case r = <-chans[0]:
		case r = <-chans[1]:
		case r = <-chans[2]:
		case <-ctx.Done():
		return
		}
		select {
		case out <- r:
		case <-ctx.Done():
			return
		}
	}
}

func merge4[T any](ctx context.Context, out Chan[T], chans []Chan[T]) {
	for {
		var r Result[T]
		select {
		case r = <-chans[0]:
		case r = <-chans[1]:
		case r = <-chans[2]:
		case r = <-chans[3]:
		case <-ctx.Done():
		return
		}
		select {
		case out <- r:
		case <-ctx.Done():
			return
		}
	}
}

func merge5[T any](ctx context.Context, out Chan[T], chans []Chan[T]) {
	for {
		var r Result[T]
		select {
		case r = <-chans[0]:
		case r = <-chans[1]:
		case r = <-chans[2]:
		case r = <-chans[3]:
		case r = <-chans[4]:
		case <-ctx.Done():
		return
		}
		select {
		case out <- r:
		case <-ctx.Done():
			return
		}
	}
}
