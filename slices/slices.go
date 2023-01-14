package slices

import (
	"golang.org/x/exp/slices"
)

func Batch[T any](s []T, size int) [][]T {
	b := [][]T{}
	for len(s) != 0 {
		n := size
		if l := len(s); l < n {
			n = l
		}
		b = append(b, s[:n:n])
		s = s[n:]
	}
	return b
}

func GroupBy[T any, K comparable](s []T, f func(T) K) map[K][]T {
	m := map[K][]T{}
	for _, v := range s {
		k := f(v)
		m[k] = append(m[k], v)
	}
	return m
}

func AppendMap[T any, U any](dst []U, src []T, f func(T) U) []U {
	dst = slices.Grow(dst, len(src))[:len(src)]
	for i, v := range src {
		dst[i] = f(v)
	}
	return dst
}

func AppendFilter[T any](dst []T, src []T, f func(T) bool) []T {
	for _, v := range src {
		if f(v) {
			dst = append(dst, v)
		}
	}
	return dst
}

func AppendDelete[T any](dst []T, src []T, i, j int) []T {
	_ = src[i:j] // bounds check
	n := len(src) - (j - i)
	dst = slices.Grow(dst, n)[:n]
	copy(dst, src[:i])
	copy(dst[i:], src[j:])
	return dst[:n]
}
