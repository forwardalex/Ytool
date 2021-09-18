package assemblyInit

import (
	"Ytool/debug"
	"Ytool/enum"
	"Ytool/log"
	"Ytool/model"
	"context"
	"gopkg.in/gomail.v2"
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
	d = gomail.NewDialer(mailConf.Host, mailConf.Port, mailConf.User, mailConf.PassWord)
	return d, nil
}

func getmailConfig() (conf *model.MailConn) {
	return debug.GetMailConf()
}
