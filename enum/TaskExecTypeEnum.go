package enum

import "sync"

// 多线程运行模式
type TaskExecTypeEnum struct {
	All Enum
	Any Enum
}

var insTaskExecType *TaskExecTypeEnum
var onceTaskExecLock sync.Mutex

func newTaskExecTypeEnum() *TaskExecTypeEnum {
	return &TaskExecTypeEnum{
		All: Enum{Id: 1, Desc: "所有任务都完成"},
		Any: Enum{Id: 2, Desc: "其中任何一个完成"},
	}
}

func GetTaskExecTypeEnum() *TaskExecTypeEnum {

	if insTaskExecType == nil {

		onceTaskExecLock.Lock()
		defer onceTaskExecLock.Unlock()

		if insTaskExecType == nil {
			insTaskExecType = newTaskExecTypeEnum()
		}

	}

	return insTaskExecType
}
