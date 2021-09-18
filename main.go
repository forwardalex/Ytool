package main

import (
	_ "Ytool/debug"
	"Ytool/layzeInit"
	"Ytool/log"
	"Ytool/mail"
	"fmt"
)

func main() {
	layzeInit.RegisterAssembly()

	fmt.Println("welcome")
	err := mail.Testmail()
	if err != nil {
		log.Error("err ", err)
	}
	fmt.Println()

}
