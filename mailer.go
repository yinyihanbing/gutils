package gutils

import (
	"fmt"
	"net/smtp"
	"strings"
	"github.com/yinyihanbing/gutils/logs"
)

// 邮件发送: goutils.Mailer.SendSingle("邮箱地址", "标题", "内容")
type mailer struct {
	MailHost     string
	MailAuthUser string
	MailAuthPass string
	MailFrom     string
	MailUser     string
}

type mailMessage struct {
	To      []string
	From    string
	Subject string
	Body    string
	User    string
	Type    string
	Massive bool
}

var Mailer = new(mailer)

func init() {
	Mailer = new(mailer)
	Mailer.MailHost = "smtp.qq.com:587"
	Mailer.MailAuthUser = "yinyihanbing@qq.com"
	Mailer.MailAuthPass = "yixiyxxukujgecci"
	Mailer.MailFrom = "yinyihanbing@qq.com"
	Mailer.MailUser = "随风飘落的叶"
}

func MailerInitialize(host, authUser, authPass, from, user string) {
	Mailer.MailHost = host
	Mailer.MailAuthUser = authUser
	Mailer.MailAuthPass = authPass
	Mailer.MailFrom = from
	Mailer.MailUser = user
}

// Create New mail message use MailFrom and MailUser
func (this *mailer) NewMailMessage(To []string, subject, body string) mailMessage {
	msg := mailMessage{
		To:      To,
		From:    this.MailFrom,
		Subject: subject,
		Body:    body,
		Type:    "html",
	}
	msg.User = this.MailUser
	return msg
}

// create mail content
func (this mailMessage) Content() string {
	// set mail type
	contentType := "text/plain; charset=UTF-8"
	if this.Type == "html" {
		contentType = "text/html; charset=UTF-8"
	}

	// create mail content
	content := "From: " + this.User + "<" + this.From +
		">\r\nSubject: " + this.Subject + "\r\nContent-Type: " + contentType + "\r\n\r\n" + this.Body
	return content
}

// Direct Send mail message
func (this *mailer) Send(to []string, subject, body string) (int, error) {
	host := strings.Split(this.MailHost, ":")

	// get message body
	msg := this.NewMailMessage(to, subject, body)
	content := msg.Content()
	auth := smtp.PlainAuth("", this.MailAuthUser, this.MailAuthPass, host[0])

	if len(msg.To) == 0 {
		return 0, fmt.Errorf("empty receive emails")
	}

	if len(msg.Body) == 0 {
		return 0, fmt.Errorf("empty email body")
	}

	if msg.Massive {
		// send mail to multiple emails one by one
		num := 0
		for _, to := range msg.To {
			body := []byte("To: " + to + "\r\n" + content)
			err := smtp.SendMail(this.MailHost, auth, msg.From, []string{to}, body)
			if err != nil {
				return num, err
			}
			num++
		}
		return num, nil
	} else {
		body := []byte("To: " + strings.Join(msg.To, ";") + "\r\n" + content)

		// send to multiple emails in one message
		err := smtp.SendMail(this.MailHost, auth, msg.From, msg.To, body)
		if err != nil {
			return 0, err
		} else {
			return 1, nil
		}
	}
}

// Async Send mail message
func (this *mailer) SendAsync(to []string, subject, body string) {
	// TODO may be need pools limit concurrent nums
	go func() {
		if num, err := this.Send(to, subject, body); err != nil {
			tos := strings.Join(to, "; ")
			logs.Error(fmt.Sprintf("Async send email %d succeed, not send emails: %s err: %s", num, tos, err))
		}
	}()
}

// 发送邮件, to=接收人, subject=标题, body=内容
func (this *mailer) SendSingle(to string, subject, body string) (int, error) {
	return this.Send([]string{to}, subject, body)
}

// 异步发送邮件, to=接收人, subject=标题, body=内容
func (this *mailer) SendAsyncSingle(to string, subject, body string) {
	this.SendAsync([]string{to}, subject, body)
}
