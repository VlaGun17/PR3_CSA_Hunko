package smtp

import (
	"fmt"
	"mime"
	"net/smtp"
	"strings"
	"time"
)

type Config struct {
	Host     string
	Port     string
	From     string
	Password string
}

type Sender struct {
	config Config
}

func NewSender(config Config) *Sender {
	return &Sender{config: config}
}

func (s *Sender) SendEmail(to []string, subject, body string) error {
	smtpAuth := smtp.PlainAuth("", s.config.From, s.config.Password, s.config.Host)
	headers := []string{
		fmt.Sprintf("From: %s", s.config.From),
		fmt.Sprintf("To: %s", strings.Join(to, ", ")),
		fmt.Sprintf("Subject: %s", mime.QEncoding.Encode("UTF-8", subject)),
		fmt.Sprintf("Date: %s", time.Now().Format(time.RFC1123Z)),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"UTF-8\"",
	}
	message := []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + body)
	return smtp.SendMail(s.config.Host+":"+s.config.Port, smtpAuth, s.config.From, to, message)
}
