package layzeInit

import (
	"Ytool/assemblyInit"
	"Ytool/enum"
	"context"
)

var assemblyMap = make(map[enum.Enum]assemblyInit.Assembly)

/*
注册组件实现
*/
func RegisterAssembly() {
	addAssembly(&assemblyInit.MySqlInit{})
	addAssembly(&assemblyInit.MailInit{})

}

func addAssembly(impl assemblyInit.Assembly) {
	assemblyMap[impl.GetAssemblyType()] = impl
}

func GetAssembly(aType enum.Enum) interface{} {
	if len(assemblyMap) == 0 {
		return nil
	}
	impl, ok := assemblyMap[aType]

	if impl == nil || !ok {
		return nil
	}

	if impl.GetAssemblyObj() != nil && impl.GetAssemblyObj().GetObj() != nil {
		return impl.GetAssemblyObj().GetObj()
	}

	background := context.Background()

	impl.GetAssemblyObj().SetObj(impl.InitAssembly(background))

	return impl.GetAssemblyObj().GetObj()
}
