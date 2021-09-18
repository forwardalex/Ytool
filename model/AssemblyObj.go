package model

type AssemblyObj struct {
	obj interface{}
}

func (a *AssemblyObj) GetObj() interface{} {
	return a.obj
}

func (a *AssemblyObj) SetObj(ob interface{}) {
	a.obj = ob
}
