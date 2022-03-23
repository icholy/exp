package chans

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func assertNilErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertErr(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func assertDeepEqual(t *testing.T, want, got any) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %#v, got %#v", want, got)
	}
}

func TestChan(t *testing.T) {
	ctx := context.Background()
	ch := Go(ctx, func() (string, error) {
		return "hello world", nil
	})
	s, err := ch.Recv(nil)
	assertNilErr(t, err)
	assertDeepEqual(t, s, "hello world")
}

func TestChanErr(t *testing.T) {
	ctx := context.Background()
	ch := Go(ctx, func() (string, error) {
		return "", errors.New("failure")
	})
	s, err := ch.Recv(nil)
	assertErr(t, err)
	assertDeepEqual(t, s, "")
}

func TestChanContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ch := Go(ctx, func() (int, error) {
		return 42, nil
	})
	x, err := ch.Recv(ctx)
	assertDeepEqual(t, x, 0)
	assertDeepEqual(t, err, context.Canceled)
}
