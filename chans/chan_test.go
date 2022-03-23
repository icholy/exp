package chans

import (
	"context"
	"errors"
	"testing"
)

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
