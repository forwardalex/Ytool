package assemblyInit

import (
	"Ytool/enum"
	"Ytool/model"
	"context"
)

type Assembly interface {
	// 组件初始化
	InitAssembly(ctx context.Context) interface{}
	// 组件类型
	GetAssemblyType() enum.Enum

	GetAssemblyObj() *model.AssemblyObj
}
