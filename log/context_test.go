package log

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeader(t *testing.T) {
	ctx := context.Background()
	ctx, err := WriteHeader(&ctx, "test", "value")
	if !assert.Equal(t, err, nil) {
		t.Error("failed")
	}
	md := GetOutHeader(ctx)
	if !assert.Equal(t, md["test"], []string{"value"}) {
		t.Error("failed")
	}
}
