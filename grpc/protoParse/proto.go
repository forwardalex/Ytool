package protoParse

import (
	"context"
	"github.com/emicklei/proto"
	"github.com/forwardalex/Ytool/log"
	"io/ioutil"
	"os"
	"strings"
)

type Service struct {
	Name      string
	RpcMethod []string
	Option    string
}

var (
	S          = make([]Service, 0) //Service 数量
	serviceStr string
	method     []string
	option     string
)

func NewProto(fileName string) {
	var parser *proto.Parser
	method = make([]string, 0)
	serviceStr = ""
	option = ""
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
		proto.WithRPC(handleRPC),
		proto.WithOption(handleOption),
	)

	S = append(S, Service{
		Name:      serviceStr,
		RpcMethod: method,
		Option:    option,
	})
	serviceStr = ""
	option = ""
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

func handleOption(o *proto.Option) {
	option = o.Constant.Source
}

const srcDir = "./grpc/protoParse/testDir/"

var (
	fileMap = make(map[string]bool, 0)
)

func ParseProtobuf() error {
	err := searchfile()
	if err != nil {
		return err
	}
	for k, v := range fileMap {
		if v == false {
			NewProto(srcDir + k)
		}
		fileMap[k] = true
	}
	return nil
	//expr, _ := parser.ParseFile()
	//fmt.Printf("%#v\n", expr)
}

func searchfile() error {
	fs, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Error(context.Background(), "read ./src dir failed  ", err)
		return err
	}
	for _, v := range fs {
		if strings.Contains(v.Name(), ".proto") {
			_, ok := fileMap[v.Name()]
			if !ok {
				fileMap[v.Name()] = false
			}
		}
	}
	return nil
}
