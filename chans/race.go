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
    var r Result[T]
	select {
	case r = <-chans[0]:
	case <-ctx.Done():
		return Result[T]{Err: ctx.Err()}
	}
    return r
}

func race2[T any](ctx context.Context, chans []Chan[T]) Result[T] {
    var r Result[T]
	select {
	case r = <-chans[0]:
	case r = <-chans[1]:
	case <-ctx.Done():
		return Result[T]{Err: ctx.Err()}
	}
    return r
}

func race3[T any](ctx context.Context, chans []Chan[T]) Result[T] {
    var r Result[T]
	select {
	case r = <-chans[0]:
	case r = <-chans[1]:
	case r = <-chans[2]:
	case <-ctx.Done():
		return Result[T]{Err: ctx.Err()}
	}
    return r
}

func race4[T any](ctx context.Context, chans []Chan[T]) Result[T] {
    var r Result[T]
	select {
	case r = <-chans[0]:
	case r = <-chans[1]:
	case r = <-chans[2]:
	case r = <-chans[3]:
	case <-ctx.Done():
		return Result[T]{Err: ctx.Err()}
	}
    return r
}

func race5[T any](ctx context.Context, chans []Chan[T]) Result[T] {
    var r Result[T]
	select {
	case r = <-chans[0]:
	case r = <-chans[1]:
	case r = <-chans[2]:
	case r = <-chans[3]:
	case r = <-chans[4]:
	case <-ctx.Done():
		return Result[T]{Err: ctx.Err()}
	}
    return r
}

func raceN[T any](ctx context.Context, chans []Chan[T]) Result[T] {
    bchans := slices.Map(
        slices.Batch(chans, 5),
        func(batch []Chan[T]) Chan[T] {
            ch := make(Chan[T])
            go func() {
                var r Result[T]
                r.Value, r.Err = Race(ctx, batch...)
                select {
                case ch <- r:
                case <-ctx.Done():
                }
            }()
            return ch
        },
    )
	var r Result[T]
	r.Value, r.Err = Race(ctx, bchans...)
	return r
}
