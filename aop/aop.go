package aop

import (
	"Ytool/log"
	"bou.ke/monkey"

	"reflect"
	"time"
)

//连接点
type JoinPoint struct {
	Receiver  interface{}
	Method    reflect.Method
	Params    []reflect.Value
	Result    []reflect.Value
	beginTime time.Time
}

func NewJoinPoint(receiver interface{}, params []reflect.Value, method reflect.Method) *JoinPoint {
	point := &JoinPoint{
		Receiver:  receiver,
		Params:    params,
		Method:    method,
		beginTime: time.Now(),
	}
	fn := method.Func
	fnType := fn.Type()
	nout := fnType.NumOut()
	point.Result = make([]reflect.Value, nout)
	for i := 0; i < nout; i++ {
		//默认返回空值
		point.Result[i] = reflect.Zero(fnType.Out(i))
	}

	return point
}

//切面接口
type AspectInterface interface {
	Before(point *JoinPoint) (bool, error)
	After(point *JoinPoint)
	Finally(point *JoinPoint)
}

//注册切点
func RegisterPoint(pointType reflect.Type, aspect AspectInterface) {
	for i := 0; i < pointType.NumMethod(); i++ {
		method := pointType.Method(i)
		//方法位置字符串"包名.接收者.方法名"，用于匹配代理
		//methodLocation := fmt.Sprintf("%s.%s.%s", pkgPth, receiverName, method.Name)
		var guard *monkey.PatchGuard
		//var patch *monkey.Patches
		var proxy = func(in []reflect.Value) []reflect.Value {
			guard.Unpatch()
			//defer patch.Reset()
			defer guard.Restore()
			receiver := in[0]
			point := NewJoinPoint(receiver, in[1:], method)
			defer finallyProcessed(point, aspect)
			if before, err := beforeProcessed(point, aspect); !before || err != nil {
				response := reflect.New(point.Result[0].Type())
				point.Result[0] = response.Elem()
				if err != nil && len(point.Result) > 1 {
					point.Result[1] = reflect.ValueOf(err)
				}
				return point.Result
			}
			point.Result = receiver.MethodByName(method.Name).Call(in[1:])
			afterProcessed(point, aspect)
			return point.Result
		}
		//动态创建代理函数   把proxy加入需要测试的函数中
		proxyFn := reflect.MakeFunc(method.Func.Type(), proxy)
		//利用monkey框架替换代理函数   替换函数
		guard = monkey.PatchInstanceMethod(pointType, method.Name, proxyFn.Interface())
		//patch =monkey.ApplyMethod(pointType, method.Name, proxyFn.Interface())
	}
}

//注册切面
//func RegisterAspect(aspect AspectInterface) {
//	aspectList = append(aspectList, aspect)
//}

//前置处理
func beforeProcessed(point *JoinPoint, aspect AspectInterface) (bool, error) {
	log.Info("before", nil)
	if &aspect != nil {
		before, err := aspect.Before(point)
		if !before || err != nil {
			return before, err
		}
	}
	return true, nil
}

//后置处理
func afterProcessed(point *JoinPoint, aspect AspectInterface) {
	log.Info("after", nil)
	if &aspect != nil {
		aspect.After(point)
	}
}

//最终处理
func finallyProcessed(point *JoinPoint, aspect AspectInterface) {
	log.Info("end", nil)
	if &aspect != nil {
		aspect.Finally(point)
	}
}
