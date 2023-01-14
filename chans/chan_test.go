package chans

import (
	"context"
	"errors"
	"testing"

	"github.com/icholy/exp/internal/assert"

	"golang.org/x/sync/errgroup"
)

func TestChan(t *testing.T) {
	var g errgroup.Group
	ch := Go(&g, func() (string, error) {
		return "hello world", nil
	})
	s, err := ch.Recv()
	assert.NoErr(t, err)
	assert.DeepEqual(t, s, "hello world")
}

func TestChanErr(t *testing.T) {
	var g errgroup.Group
	ch := Go(&g, func() (string, error) {
		return "", errors.New("failure")
	})
	s, err := ch.Recv()
	assert.Err(t, err)
	assert.DeepEqual(t, s, "")
}

func TestChanContext(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	g, ctx := errgroup.WithContext(ctx)
	ch := Go(g, func() (int, error) {
		return 42, nil
	})
	x, err := Recv(ctx, ch)
	assert.DeepEqual(t, 0, x)
	assert.DeepEqual(t, err, context.Canceled)
}
