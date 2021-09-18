package main

import (
	"Ytool/log"
	"Ytool/mail"
	"Ytool/tool"
	"fmt"
)

func main() {
	tool.Init()
	fmt.Println("welcome")
	err := mail.Testmail()
	if err != nil {
		log.Error("err ", err)
	}
	fmt.Println()

}
