package main

import (
	_ "Ytool/debug"
	"Ytool/test"
	"fmt"
)

func main() {
	fmt.Println("welcome")
	test.TestBlame()
	fmt.Println("over")

}
