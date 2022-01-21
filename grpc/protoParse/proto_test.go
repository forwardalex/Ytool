package protoParse

import (
	"context"
	"fmt"
	"github.com/forwardalex/Ytool/log"
	"io/ioutil"
	"strings"
	"testing"
)

var (
	fileMap = make(map[string]bool, 0)
)

const srcDir = "./testDir/"

func searchfile() {
	fs, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Error(context.Background(), "read ./src dir failed  ", err)
	}
	for _, v := range fs {
		if strings.Contains(v.Name(), ".proto") {
			_, ok := fileMap[v.Name()]
			if !ok {
				fileMap[v.Name()] = false
			}
		}
	}
}
func ParseProtobuf() {
	searchfile()
	for k, v := range fileMap {
		if v == false {
			NewProto(srcDir + k)
		}
		fileMap[k] = true
	}
	//expr, _ := parser.ParseFile()
	//fmt.Printf("%#v\n", expr)
}

func TestNewProto(t *testing.T) {
	ParseProtobuf()
	for k, v := range S {
		fmt.Println("service", k, v.Name)
		for _, method := range v.RpcMethod {
			fmt.Println("method ", method)
		}
	}
}
