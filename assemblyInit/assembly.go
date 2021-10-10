package assemblyInit

import (
	"context"
	"github.com/forwardalex/Ytool/enum"
	"github.com/forwardalex/Ytool/model"
)

type Assembly interface {
	// 组件初始化
	InitAssembly(ctx context.Context) interface{}
	// 组件类型
	GetAssemblyType() enum.Enum

	GetAssemblyObj() *model.AssemblyObj
}
