package protoParse

import (
	"context"
	"fmt"
	"github.com/emicklei/proto"
	"github.com/forwardalex/Ytool/log"
	"os"
)

var (
	Service = &proto.Service{}
	RPC     = &proto.RPC{}
)

func NewProto(fileName string) {
	var parser *proto.Parser
	file, err := os.Open(fileName)
	if err != nil {
		log.Error(context.Background(), "open file failed ", err)
		return
	}
	defer file.Close()
	parser = proto.NewParser(file)
	definition, _ := parser.Parse()

	proto.Walk(definition,
		proto.WithService(handleService),
		proto.WithRPC(handleRPC),
		proto.WithMessage(handleMessage))
}
func handleService(s *proto.Service) {
	Service = s
}

func handleMessage(m *proto.Message) {
	fmt.Println(m.Name)
	fmt.Println(m.Elements)
}
func handleRPC(r *proto.RPC) {
	RPC = r
	fmt.Println(r.Name)
}
