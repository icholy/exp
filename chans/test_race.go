package chans

import (
	"context"
	"testing"

	"github.com/icholy/exp/internal/assert"
)

func TestRace1(t *testing.T) {
	ch := Go(nil, func() (string, error) {
		return "test", nil
	})
	ctx := context.Background()
	v, err := Race(ctx, ch)
	assert.NoErr(t, err)
	assert.DeepEqual(t, "test", v)
}
