package log

import (
	"context"
	"testing"
)

func TestBlame(t *testing.T) {
	c, err := FindCommit(context.Background(), "../main.go", 1, nil)
	if err != nil {
		Error("test ", err)
	}
	if c.CommitName == "" {
		panic("err")
	}
}
