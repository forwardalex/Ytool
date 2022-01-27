package log

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

	ctx = context.Background()
	traceId := uuid.New().String()
	fmt.Println("tarce ", traceId)
	ctx, err = WriteHeader(&ctx, TraceStringKey, traceId)
	if err != nil {
		fmt.Println("err ", err)
	}
	md = GetOutHeader(ctx)
	fmt.Println(md)
	fmt.Println("====11111", md[TraceStringKey])
}
