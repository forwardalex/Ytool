package recover

//RecoverFromPanic panic recover
//defer func
//func RecoverFromPanic(funcName string) {
//	if e := recover(); e != nil {
//		buf := make([]byte, 64<<10)
//		buf = buf[:runtime.Stack(buf, false)]
//
//		log.Errorf("[%s] func_name: %v, stack: %s", funcName, e, string(buf))
//
//		panicError := fmt.Errorf("%v", e)
//		panic_reporter_client.ReportPanic(panicError.Error(), funcName, string(buf))
//	}
//
//	return
//}
