package recover

import (
	"Ytool/env"
	"Ytool/log"
	"Ytool/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var (
	panicReportOnce sync.Once
	//todo 搞个服务
	url = "xxx.url"
)

//RecoverFromPanic panic recover
//defer func

func RecoverFromPanic(funcName string) {
	if e := recover(); e != nil {
		buf := make([]byte, 64<<10)
		buf = buf[:runtime.Stack(buf, false)]

		log.Errorf("[%s] func_name: %v, stack: %s", funcName, e, string(buf))

		panicError := fmt.Errorf("%v", e)
		ReportPanic(panicError.Error(), funcName, string(buf))
	}
	return
}

func ReportPanic(errInfo, funcName, stack string) (err error) {
	panicReportOnce.Do(func() {
		defer func() { recover() }()
		go func() {
			panicReq := &model.PanicReq{
				Service:   env.Service(),
				ErrorInfo: errInfo,
				Stack:     stack,
				FuncName:  funcName,
				Host:      env.HostIP(),
				PodName:   env.PodName(),
			}
			var jsonBytes []byte
			jsonBytes, err = json.Marshal(panicReq)
			if err != nil {
				return
			}
			var req *http.Request
			req, err = http.NewRequest("GET", url, bytes.NewBuffer(jsonBytes))
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 5 * time.Second}
			var resp *http.Response
			resp, err = client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			return
		}()
	})

	return
}
