package mail

import (
	"Ytool/debug"
	"Ytool/enum"
	"Ytool/layzeInit"
	"gopkg.in/gomail.v2"
)

var MailDail *gomail.Dialer

func GetMailConn() *gomail.Dialer {
	if MailDail != nil {
		return MailDail
	}
	assembly := layzeInit.GetAssembly(enum.GetAssemblyEnum().Mail)

	if assembly == nil {
		return nil
	}

	return assembly.(*gomail.Dialer)
}

func SendMail(mailTo []string, subject string, body string, d *gomail.Dialer) error {
	conf := debug.GetMailConf()
	m := gomail.NewMessage()
	m.SetHeader("From", "panic info "+"<"+conf.User+">") //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo...)                         //发送给多个用户
	m.SetHeader("Subject", subject)                      //设置邮件主题
	m.SetBody("text/html", body)                         //设置邮件正文
	err := d.DialAndSend(m)
	return err
}
func Testmail() error {
	//定义收件人
	mailTo := []string{
		"550507808@qq.com",
	}
	//邮件主题为"Hello"
	subject := "server panic"
	// 邮件正文
	body := "this is a test mail"
	err := SendMail(mailTo, subject, body, GetMailConn())
	if err != nil {
		return err
	}
	return nil
}
