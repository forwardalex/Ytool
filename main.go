package main

import (
	"fmt"
	"github.com/forwardalex/Ytool/grpcInterceptor"
	"github.com/forwardalex/Ytool/tool"
)

func main() {
	tool.Init()
	fmt.Println("welcome")
	//test.TestBlame()
	grpcInterceptor.Test()

}
