package slices

func Batch[T any](s []T, size int) [][]T {
	b := [][]T{}
	for len(s) != 0 {
		n := size
		if l := len(s); l < n {
			n = l
		}
		b = append(b, s[:n])
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

func InPlaceFilter[T any](s []T, f func(T) bool) []T {
	i := 0
	for _, v := range s {
		if f(v) {
			s[i] = v
			i++
		}
	}
	return s[:i]
}

func AppendFilter[T any](dst []T, src []T,  f func(T) bool) []T {
	for _, v := range src {
		if f(v) {
			dst = append(dst, v)
		}
	}
	return dst
}

func Filter[T any](s []T, f func(T) bool) []T {
	return AppendFilter(nil, s, f)
}
