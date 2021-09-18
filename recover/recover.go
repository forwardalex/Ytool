package rec

import (
	"Ytool/env"
	"Ytool/log"
	"Ytool/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
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
		ctx := context.Background()
		buf := make([]byte, 64<<10)
		buf = buf[:runtime.Stack(buf, false)]
		_, fn, line, _ := runtime.Caller(2)
		commit, err := log.FindCommit(ctx, fn, line, nil)
		if err != nil {
			log.Error("error ", err.Error())
		}
		panicError := fmt.Errorf("%v", e)
		getenv := os.Getenv("ENV_NAME")
		log.Infof("commit is %s, panic err is %s env is %s ", commit, panicError, getenv)
		switch getenv {
		case env.EnvMapStr[env.EnvLocal], env.EnvMapStr[env.EnvDevelopment]:
			log.Infof("panic commit is %s,commit filename is %s line is %d ", commit.CommitName, fn, line)
			err = ReportPanicCommit(panicError.Error(), funcName, fn, string(buf), line, commit)
			if err != nil {
				log.Error("report failed ", err)
			}
		case env.EnvMapStr[env.EnvTest], env.EnvMapStr[env.EnvPreRelease], env.EnvMapStr[env.EnvProduction]:
			err = ReportPanicCommit(panicError.Error(), funcName, fn, string(buf), line, commit)
			if err != nil {
				log.Error("report failed ", err)
			}
		}
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

func ReportPanicCommit(errInfo, funcName, stack, fileName string, line int, commit log.BlameLine) (err error) {
	panicReportOnce.Do(func() {
		mailBody := fmt.Sprintf(`
	Service  	: %s,
	ErrorInfo 	: %s,
	Stack 		: %s,
	Host		: %s,
	funcName 	: %s,
	PodName     : %s,
	LastCommitUser		: %s,
	CommitTime			: %s,
	FileName			: %s,
	Line                : %d, `, env.Service(), errInfo, stack, env.HostIP(), funcName, env.PodName(), commit.CommitName, commit.CommitDate.String(), fileName, line)
		handleUser := strings.Trim(commit.CommitEmail, "<>")
		fmt.Println(handleUser)
		fmt.Println(mailBody)
		//mailTo:=[]string{handleUser}
		//err=mail.SendMail(mailTo,"server panic",mailBody,mail.GetMailConn())
		if err != nil {
			log.Error("send mail failed ", err)
		}
	})
	return nil
}
