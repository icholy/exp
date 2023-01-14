package chans

import (
	"sync"

	"github.com/icholy/exp/slices"
)

func Merge[T any](chans ...chan T) chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		merge(ch, chans)
	}()
	return ch
}

func without[T any](chans []chan T, i int) []chan T {
	return slices.AppendDelete(nil, chans, i, i+1)
}

func merge[T any](out chan T, chans []chan T) {
	switch len(chans) {
	case 1:
		merge1(out, chans)
	case 2:
		merge2(out, chans)
	case 3:
		merge3(out, chans)
	case 4:
		merge4(out, chans)
	case 5:
		merge5(out, chans)
	default:
		var wg sync.WaitGroup
		for _, batch := range slices.Batch(chans, 5) {
			wg.Add(1)
			batch := batch
			go func() {
				defer wg.Done()
				merge(out, batch)
			}()
		}
		wg.Wait()
	}
}

func merge1[T any](out chan T, chans []chan T) {
	for {
		r, ok := <-chans[0]
		if !ok {
			return
		}
		out <- r
	}
}

func merge2[T any](out chan T, chans []chan T) {
	for {
		var i int
		var r T
		var ok bool
		select {
		case r, ok = <-chans[0]:
			i = 0
		case r, ok = <-chans[1]:
			i = 1
		}
		if !ok {
			merge1(out, without(chans, i))
			return
		}
		out <- r
	}
}

func merge3[T any](out chan T, chans []chan T) {
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
		}
		if !ok {
			merge2(out, without(chans, i))
			return
		}
		out <- r
	}
}

func merge4[T any](out chan T, chans []chan T) {
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
		}
		if !ok {
			merge3(out, without(chans, i))
			return
		}
		out <- r
	}
}

func merge5[T any](out chan T, chans []chan T) {
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
		}
		if !ok {
			merge4(out, without(chans, i))
			return
		}
		out <- r
	}
}
