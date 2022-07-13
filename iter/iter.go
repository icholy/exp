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
