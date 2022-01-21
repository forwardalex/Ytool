package protoParse

import (
	"context"
	"github.com/emicklei/proto"
	"github.com/forwardalex/Ytool/log"
	"os"
)

type Service struct {
	Name      string
	RpcMethod []string
}

var (
	S          = make([]Service, 0) //Service 数量
	serviceStr string
	method     []string
)

func NewProto(fileName string) {
	var parser *proto.Parser
	method = make([]string, 0)
	serviceStr = ""
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
		proto.WithMessage(handleMessage),
		proto.WithRPC(handleRPC))

	S = append(S, Service{
		Name:      serviceStr,
		RpcMethod: method,
	})
	serviceStr = ""
	method = []string{}
}

func handleService(s *proto.Service) {
	serviceStr = s.Name
	log.Info(context.Background(), "service Name", s.Name)
}

func handleMessage(m *proto.Message) {
	//fmt.Println(m.Name)
	//fmt.Println(m.Comment)
}
func handleRPC(r *proto.RPC) {
	method = append(method, r.Name)
	log.Info(context.Background(), "rpc Name", r.Name)
}
