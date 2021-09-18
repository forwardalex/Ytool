package enum

import "sync"

type AssemblyEnum struct {
	Kafka        Enum
	MySQL        Enum
	PrivateCloud Enum
	RedLock      Enum
	TxMeeting    Enum
	Redis        Enum
	Mail         Enum
}

var insAssembly *AssemblyEnum
var onceAssembly sync.Once

func newAssemblyEnum() *AssemblyEnum {
	return &AssemblyEnum{
		Kafka:        Enum{Id: 1, Desc: ""},
		MySQL:        Enum{Id: 2, Desc: ""},
		PrivateCloud: Enum{Id: 3, Desc: ""},
		RedLock:      Enum{Id: 4, Desc: ""},
		TxMeeting:    Enum{Id: 5, Desc: ""},
		Redis:        Enum{Id: 6, Desc: ""},
		Mail:         Enum{Id: 7, Desc: ""},
	}
}

func GetAssemblyEnum() *AssemblyEnum {

	if insAssembly == nil {
		onceAssembly.Do(func() {
			if insAssembly == nil {
				insAssembly = newAssemblyEnum()
			}
		})
	}

	return insAssembly
}
