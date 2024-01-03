package utils

import (
	conf "funoj-backend/config"
	"gopkg.in/gomail.v2"
)

type SysEmailMessage struct {
	To      []string
	Subject string
	Body    string
}

func SendSysEmail(config *conf.EmailConfig, message SysEmailMessage) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(config.User, "funoj")) // 添加别名
	m.SetHeader("To", message.To...)                           // 发送给用户(可以多个)
	m.SetHeader("Subject", message.Subject)                    // 设置邮件主题
	m.SetBody("text/html", message.Body)                       // 设置邮件正文
	d := gomail.NewDialer(config.Host, config.Port, config.User, config.Password)
	err := d.DialAndSend(m)
	return err
}
