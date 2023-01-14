package chans

import (
	"context"
	"testing"

	"github.com/icholy/exp/internal/assert"
	"golang.org/x/exp/slices"
)

func TestMerge1(t *testing.T) {
	m := Merge(
		context.Background(),
		sliceToChan([]string{"a", "b", "c"}),
	)
	assert.DeepEqual(t, []string{"a", "b", "c"}, chanToSlice(m))
}

func TestMerge2(t *testing.T) {
	m := Merge(
		context.Background(),
		sliceToChan([]int{1, 2, 3}),
		sliceToChan([]int{4, 5, 6}),
	)
	s := chanToSlice(m)
	slices.Sort(s)
	assert.DeepEqual(t, s, []int{1, 2, 3, 4, 5, 6})
}

func sliceToChan[T any](s []T) chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, v := range s {
			ch <- v
		}
	}()
	return ch
}

func chanToSlice[T any](ch chan T) []T {
	var s []T
	for v := range ch {
		s = append(s, v)
	}
	return s
}
