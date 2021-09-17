package main

import (
	_ "Ytool/debug"
	"Ytool/log"
	"bytes"
	"context"
	"fmt"
)

var stderr bytes.Buffer

func main() {
	ctx := context.Background()
	err := log.FindCommit(ctx, "./work-config.yml", 10, nil)
	if err != nil {
		fmt.Println(err)
	}
}
