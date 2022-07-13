package iter

import (
	"reflect"
	"testing"
)

func TestBatchIter(t *testing.T) {

	var idx int
	batches := [][]string{
		{"a", "b"},
		{"b"},
		{"c", "d"},
	}

	it := &BatchIter[string]{
		NextBatch: func() ([]string, bool, error) {
			batch := batches[idx]
			idx++
			return batch, idx < len(batches), nil
		},
	}
	want := []string{"a", "b", "b", "c", "d"}
	got, err := Slice[string](it)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v, got %s", want, got)
	}
}
