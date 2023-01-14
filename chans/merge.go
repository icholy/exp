package chans

import (
	"context"
	"sync"

	"github.com/icholy/exp/slices"
)

func Merge[T any](ctx context.Context, chans ...chan T) chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		merge(ctx, ch, chans)
	}()
	return ch
}

func merge[T any](ctx context.Context, out chan T, chans []chan T) {
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
		var wg sync.WaitGroup
		for _, batch := range slices.Batch(chans, 5) {
			wg.Add(1)
			batch := batch
			go func() {
				defer wg.Done()
				merge(ctx, out, batch)
			}()
		}
		wg.Wait()
	}
}

func merge1[T any](ctx context.Context, out chan T, chans []chan T) {
	for {
		var r T
		var ok bool
		select {
		case r, ok = <-chans[0]:
		case <-ctx.Done():
			return
		}
		if !ok {
			return
		}
		select {
		case out <- r:
		case <-ctx.Done():
			return
		}
	}
}

func merge2[T any](ctx context.Context, out chan T, chans []chan T) {
	for {
		var i int
		var r T
		var ok bool
		select {
		case r, ok = <-chans[0]:
			i = 0
		case r, ok = <-chans[1]:
			i = 1
		case <-ctx.Done():
			return
		}
		if !ok {
			merge1(ctx, out, without(chans, i))
			return
		}
		select {
		case out <- r:
		case <-ctx.Done():
			return
		}
	}
}

func without[T any](chans []chan T, i int) []chan T {
	return slices.AppendDelete(chans, chans, i, i+1)
}

func merge3[T any](ctx context.Context, out chan T, chans []chan T) {
	for {
		var i int
		var r T
		var ok bool
		select {
		case r, ok = <-chans[0]:
			i = 0
		case r, ok = <-chans[1]:
			i = 1
		case r, ok = <-chans[2]:
			i = 2
		case <-ctx.Done():
			return
		}
		if !ok {
			merge2(ctx, out, without(chans, i))
			return
		}
		select {
		case out <- r:
		case <-ctx.Done():
			return
		}
	}
}

func merge4[T any](ctx context.Context, out chan T, chans []chan T) {
	for {
		var i int
		var r T
		var ok bool
		select {
		case r, ok = <-chans[0]:
			i = 0
		case r, ok = <-chans[1]:
			i = 1
		case r, ok = <-chans[2]:
			i = 2
		case r, ok = <-chans[3]:
			i = 3
		case <-ctx.Done():
			return
		}
		if !ok {
			merge3(ctx, out, without(chans, i))
			return
		}
		select {
		case out <- r:
		case <-ctx.Done():
			return
		}
	}
}

func merge5[T any](ctx context.Context, out chan T, chans []chan T) {
	for {
		var i int
		var r T
		var ok bool
		select {
		case r, ok = <-chans[0]:
		case r, ok = <-chans[1]:
		case r, ok = <-chans[2]:
		case r, ok = <-chans[3]:
		case r, ok = <-chans[4]:
		case <-ctx.Done():
			return
		}
		if !ok {
			merge4(ctx, out, without(chans, i))
			return
		}
		select {
		case out <- r:
		case <-ctx.Done():
			return
		}
	}
}
