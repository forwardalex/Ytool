package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/forwardalex/Ytool/grpc/protoParse"
	"github.com/forwardalex/Ytool/log"
	"os"
)

func main() {
	Init()
	fmt.Println("ok")
	Genfile(context.Background())

}

var methodStr = `
//%s 
func (server *%s) %s(ctx context.Context,
	req *pb.%sReq) (resp *pb.%sResp, err error) {
	url, err := dao.%s(ctx, req)
	if err != nil {
		log.Error(ctx, "%s failed ", err)
		return nil, err
	}
	resp = &pb.%sResp{}
	return resp, nil
}
`

var packImport = `package service
import (
	"context"
	pb "%s"
)

type %s struct {
 pb.Unimplemented%sServer
}
`
var filepath = "./grpc/protoParse/gen/"

func Init() {
	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			log.Error(nil, "err ", err)
			err = os.MkdirAll(filepath, os.ModePerm)
			if err != nil {
				log.Error(nil, "err ", err)
			}
		}
	}
}

func Genfile(ctx context.Context) {
	err := protoParse.ParseProtobuf()
	if err != nil {
		log.Error(ctx, " ParseProtobuf() failed  ", err)
	}
	for k, v := range protoParse.S {
		filename := filepath + v.Name + ".go"
		fmt.Println(filename)
		fmt.Println("service", k, v.Name)
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Error(ctx, "create file failed ", err)
			return
		}
		writer := bufio.NewWriter(f)
		improtwri := packImport
		improtwri = fmt.Sprintf(improtwri, v.Option, v.Name, v.Name)
		writer.WriteString(improtwri)
		for _, method := range v.RpcMethod {
			fmt.Println("method ", method)
			methodwri := methodStr
			methodwri = fmt.Sprintf(methodwri, method, v.Name, method, method, method, method, method, method)
			writer.WriteString(methodwri)
		}
		writer.Flush() //将缓存中的内容写入文件
		f.Close()
	}

}
