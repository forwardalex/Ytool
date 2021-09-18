package main

import (
	"Ytool/test"
	"Ytool/tool"
	"fmt"
)

func main() {
	tool.Init()
	fmt.Println("welcome")
	test.TestBlame()
	fmt.Println("over")

}
