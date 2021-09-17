package main

import (
	_ "Ytool/Debug"
	"Ytool/log"
	//logs "github.com/sirupsen/logrus"
)

func main() {
	log.Infof("key is %d , another %d", 10, 10)
	//fmt.Printf("key is %d , another %d\n",10,10)
	//fmt.Printf("%d,%s",10,"s")
}
