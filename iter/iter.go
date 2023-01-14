package iter

type Iter[T any] interface {
	Next() bool
	Value() T
	Err() error
}

func Slice[T any](it Iter[T]) ([]T, error) {
	ss := []T{}
	for it.Next() {
		ss = append(ss, it.Value())
	}
	return ss, it.Err()
}

type sliceIter[T any] struct {
	v T
	i int
	s []T
}

func (_ *sliceIter[T]) Err() error { return nil }

func (s *sliceIter[T]) Value() T { return s.v }

func (s *sliceIter[T]) Next() bool {
	if s.i >= len(s.s) {
		return false
	}
	s.v = s.s[s.i]
	s.i++
	return true
}

func FromSlice[T any](s []T) Iter[T] {
	return &sliceIter[T]{s: s}
}
