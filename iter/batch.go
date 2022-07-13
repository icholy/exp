package iter

import "errors"

type BatchIter[T any] struct {
	batch []T
	idx   int
	val   T
	err   error
	done  bool

	NextBatch func() ([]T, bool, error)
}

func (it *BatchIter[T]) Next() bool {
	if it.err != nil {
		return false
	}
	if n := len(it.batch); n == 0 || it.idx >= n {
		if it.done {
			return false
		}
		var more bool
		it.batch, more, it.err = it.NextBatch()
		if it.err != nil {
			return false
		}
		if more && len(it.batch) == 0 {
			it.err = errors.New("got empty batch")
			return false
		}
		it.done = !more
		it.idx = 0
	}
	it.val = it.batch[it.idx]
	it.idx++
	return true
}

func (t *BatchIter[T]) Value() T {
	return t.val
}

func (it *BatchIter[T]) Err() error {
	return it.err
}
