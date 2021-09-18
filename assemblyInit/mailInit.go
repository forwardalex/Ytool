package assemblyInit

import (
	"Ytool/debug"
	"Ytool/enum"
	"Ytool/log"
	"Ytool/model"
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type MailInit struct {
	obj model.AssemblyObj
}

func (impl *MailInit) InitAssembly(ctx context.Context) interface{} {
	conn, err := initMail()
	if err != nil {
		log.Fatal("", err.Error())
	}
	return conn
}

func (*MailInit) GetAssemblyType() enum.Enum {
	return enum.GetAssemblyEnum().Mail
}

func (impl *MailInit) GetAssemblyObj() *model.AssemblyObj {
	return &impl.obj
}

//初始化 mail
func initMail() (d *gomail.Dialer, err error) {

	mailConf := getmailConfig()

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	fmt.Println(path[:index])

	// 如果是开发环境，手动指定数据库地址

	if err != nil {
		log.Error("err ", err)
		return nil, err
	}
	d = gomail.NewDialer(mailConf.Host, mailConf.Port, mailConf.User, mailConf.PassWord)
	return d, nil
}

func getmailConfig() (conf *model.MailConn) {
	return debug.GetMailConf()
}
