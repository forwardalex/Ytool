package tool

import (
	_ "Ytool/debug"
	"Ytool/layzeInit"
)

func Init() {
	//常用工具加载
	layzeInit.RegisterAssembly()
}
