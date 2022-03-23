package chans

import (
	"context"

    "github.com/icholy/exp/slices"
)

func Race[T any](ctx context.Context, chans ...Chan[T]) (T, error) {
	var r Result[T]
	switch len(chans) {
	case 1:
		r = race1(ctx, chans)
	case 2:
		r = race2(ctx, chans)
	case 3:
		r = race3(ctx, chans)
	case 4:
		r = race4(ctx, chans)
	case 5:
		r = race5(ctx, chans)
	default:
		r = raceN(ctx, chans)
	}
	return r.Value, r.Err
}

func race1[T any](ctx context.Context, chans []Chan[T]) Result[T] {
	select {
	case r := <-chans[0]:
		return r
	case <-ctx.Done():
		return Result[T]{Err: ctx.Err()}
	}
}

func race2[T any](ctx context.Context, chans []Chan[T]) Result[T] {
	select {
	case r := <-chans[0]:
		return r
	case r := <-chans[1]:
		return r
	case <-ctx.Done():
		return Result[T]{Err: ctx.Err()}
	}
}

func race3[T any](ctx context.Context, chans []Chan[T]) Result[T] {
	select {
	case r := <-chans[0]:
		return r
	case r := <-chans[1]:
		return r
	case r := <-chans[2]:
		return r
	case <-ctx.Done():
		return Result[T]{Err: ctx.Err()}
	}
}

func race4[T any](ctx context.Context, chans []Chan[T]) Result[T] {
	select {
	case r := <-chans[0]:
		return r
	case r := <-chans[1]:
		return r
	case r := <-chans[2]:
		return r
	case r := <-chans[3]:
		return r
	case <-ctx.Done():
		return Result[T]{Err: ctx.Err()}
	}
}

func race5[T any](ctx context.Context, chans []Chan[T]) Result[T] {
	select {
	case r := <-chans[0]:
		return r
	case r := <-chans[1]:
		return r
	case r := <-chans[2]:
		return r
	case r := <-chans[3]:
		return r
	case r := <-chans[4]:
		return r
	case <-ctx.Done():
		return Result[T]{Err: ctx.Err()}
	}
}

func raceN[T any](ctx context.Context, chans []Chan[T]) Result[T] {
	batches := slices.Batch(chans, 5)
	batchchans := make([]Chan[T], len(batches))
	for i, b := range batches {
        batchchans[i] = make(Chan[T])
		go func(batch []Chan[T], ch Chan[T]) {
			var r Result[T]
			r.Value, r.Err = Race(ctx, batch...)
			select {
			case ch <- r:
			case <-ctx.Done():
			}
		}(b, batchchans[i])
	}
	var r Result[T]
	r.Value, r.Err = Race(ctx, batchchans...)
	return r
}
