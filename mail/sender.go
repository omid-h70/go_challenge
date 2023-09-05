package mail

import (
	"fmt"
	"github.com/jordan-wright/email"
	_ "github.com/jordan-wright/email"
	"net/smtp"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(subject, content string, cc, to, bcc, attachedFiles []string) error
}

type GmailSender struct {
	name          string
	fromEmailAddr string
	fromEmailPass string
}

func NewGmailSender(name, fromEmailAddr, fromEmailPass string) EmailSender {
	return &GmailSender{
		name:          name,
		fromEmailAddr: fromEmailAddr,
		fromEmailPass: fromEmailPass,
	}
}

func (sender *GmailSender) SendEmail(subject, content string, cc, to, bcc, attachedFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddr)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachedFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}
	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddr, sender.fromEmailPass, smtpAuthAddress)
	e.Send(smtpServerAddress, smtpAuth)

	return nil
}
