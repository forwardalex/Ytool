package protoParse

import (
	"fmt"
	"testing"
)

func TestNewProto(t *testing.T) {
	ParseProtobuf()
	for k, v := range S {
		fmt.Println("service", k, v.Name)
		for _, method := range v.RpcMethod {
			fmt.Println("method ", method)
		}
	}
}
