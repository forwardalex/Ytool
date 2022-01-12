package assemblyInit

import (
	"context"
	"github.com/forwardalex/Ytool/debug"
	"github.com/forwardalex/Ytool/enum"
	"github.com/forwardalex/Ytool/log"
	"github.com/forwardalex/Ytool/model"
	"gopkg.in/gomail.v2"
)

type MailInit struct {
	obj model.AssemblyObj
}

func (impl *MailInit) InitAssembly(ctx context.Context) interface{} {
	conn, err := initMail()
	if err != nil {
		log.Fatal(context.Background(), "", err.Error())
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

func getmailConfig() (conf *model.MailConf) {
	return debug.GetMailConf()
}
