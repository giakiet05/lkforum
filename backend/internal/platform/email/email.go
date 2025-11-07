package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/giakiet05/lkforum/internal/config"
)

// Sender defines the interface for an email sender.
type Sender interface {
	Send(to []string, subject, body string) error
}

// smtpSender implements the Sender interface using SMTP.
type smtpSender struct{}

// NewSMTPSender creates a new SMTP email sender.
func NewSMTPSender() Sender {
	return &smtpSender{}
}

// Send sends an email using the configured SMTP server.
func (s *smtpSender) Send(to []string, subject, body string) error {
	smtpCfg := config.Cfg.SMTP

	auth := smtp.PlainAuth("", smtpCfg.Username, smtpCfg.Password, smtpCfg.Host)

	// The message needs to be in a specific format with headers.
	// The \r\n is required by the SMTP standard.
	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		smtpCfg.Username, // From address is usually the same as the auth username
		strings.Join(to, ","),
		subject,
		body,
	))

	addr := fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port)

	return smtp.SendMail(addr, auth, smtpCfg.Username, to, msg)
}
