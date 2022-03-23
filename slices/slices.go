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
