package email

import (
	"gopkg.in/mail.v2"

	conf "github.com/YasinDoyle/e-mall/config"
)

type EmailSender struct {
	SmtpHost      string `json:"smtp_host"`
	SmtpPort      int    `json:"smtp_port"`
	SmtpEmailFrom string `json:"smtp_email_from"`
	SmtpPass      string `json:"smtp_pass"`
}

func NewEmailSender() *EmailSender {
	eConfig := conf.Config.Email
	return &EmailSender{
		SmtpHost:      eConfig.SmtpHost,
		SmtpPort:      eConfig.SmtpPort,
		SmtpEmailFrom: eConfig.SmtpEmail,
		SmtpPass:      eConfig.SmtpPass,
	}
}

// Send 发送邮件
func (s *EmailSender) Send(data, emailTo, subject string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.SmtpEmailFrom)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", data)
	port := s.SmtpPort
	if port == 0 {
		port = 465
	}
	d := mail.NewDialer(s.SmtpHost, port, s.SmtpEmailFrom, s.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
