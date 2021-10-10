package tool

import (
	_ "github.com/forwardalex/Ytool/debug"
	"github.com/forwardalex/Ytool/layzeInit"
)

func Init() {
	//常用工具加载
	layzeInit.RegisterAssembly()
}
