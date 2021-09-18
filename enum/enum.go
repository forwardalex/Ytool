package enum

/*
做为枚举类型的基类
*/
type Enum struct {
	Id   int
	Desc string
}

func (e Enum) GetId() int {
	return e.Id
}

func (e Enum) GetDesc() string {
	return e.Desc
}
