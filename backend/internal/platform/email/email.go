package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	"github.com/giakiet05/lkforum/internal/config"
)

// Sender defines the interface for an email sender.
type Sender interface {
	SendVerificationEmail(to, otp string) error
	SendPasswordResetEmail(to, otp string) error
}

// SMTPSender is an implementation of Sender that uses SMTP.
type SMTPSender struct {
	from string
	auth smtp.Auth
	addr string
}

// NewSMTPSender creates a new SMTPSender.
// If SMTP configuration is incomplete, it returns a noopSender that only logs the OTP.
func NewSMTPSender() Sender {
	smtpCfg := config.Cfg.SMTP
	if smtpCfg.User == "" || smtpCfg.Pass == "" || smtpCfg.Host == "smtp.example.com" {
		log.Println("WARNING: SMTP is not fully configured. Email sending is disabled and will be logged to console instead.")
		return &noopSender{}
	}

	auth := smtp.PlainAuth("", smtpCfg.User, smtpCfg.Pass, smtpCfg.Host)
	addr := fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port)
	// The 'from' address must be the same as the user used for authentication for many SMTP servers.
	from := smtpCfg.User

	return &SMTPSender{
		from: from,
		auth: auth,
		addr: addr,
	}
}

// SendVerificationEmail sends an email with the OTP code.
func (s *SMTPSender) SendVerificationEmail(to, otp string) error {
	// Constructing the email headers and body
	subject := "Subject: Your Verification Code for LKForum\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	data := struct {
		OTP        string
		SenderName string
	}{
		OTP:        otp,
		SenderName: config.Cfg.SMTP.SenderName,
	}

	// Using an HTML template
	t, err := template.New("email").Parse(verificationEmailTemplate)
	if err != nil {
		log.Printf("Error parsing email template: %v", err)
		return err
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		log.Printf("Error executing email template: %v", err)
		return err
	}

	// The message must be in the format of `To`, `Subject`, then the body.
	// The `from` address is passed to SendMail directly.
	headers := fmt.Sprintf("To: %s\r\n%s", to, subject)
	msg := []byte(headers + mime + body.String())

	err = smtp.SendMail(s.addr, s.auth, s.from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Verification email sent to %s", to)
	return nil
}

// SendPasswordResetEmail sends an email with the OTP code for password reset.
func (s *SMTPSender) SendPasswordResetEmail(to, otp string) error {
	subject := "Subject: Password Reset Code for LKForum\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	data := struct {
		OTP        string
		SenderName string
	}{
		OTP:        otp,
		SenderName: config.Cfg.SMTP.SenderName,
	}

	t, err := template.New("email").Parse(passwordResetEmailTemplate)
	if err != nil {
		log.Printf("Error parsing password reset email template: %v", err)
		return err
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		log.Printf("Error executing password reset email template: %v", err)
		return err
	}

	headers := fmt.Sprintf("To: %s\r\n%s", to, subject)
	msg := []byte(headers + mime + body.String())

	err = smtp.SendMail(s.addr, s.auth, s.from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send password reset email to %s: %v", to, err)
		return err
	}

	log.Printf("Password reset email sent to %s", to)
	return nil
}

// noopSender is a sender that does nothing but log. Used when SMTP is not configured.
type noopSender struct{}

func (s *noopSender) SendVerificationEmail(to, otp string) error {
	log.Printf("Email sending is disabled. Verification OTP for %s: %s", to, otp)
	return nil
}

func (s *noopSender) SendPasswordResetEmail(to, otp string) error {
	log.Printf("Email sending is disabled. Password Reset OTP for %s: %s", to, otp)
	return nil
}

const verificationEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
<style>
  .container { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 20px auto; border: 1px solid #ddd; border-radius: 5px; }
  .header { background-color: #f7f7f7; padding: 15px; text-align: center; border-bottom: 1px solid #ddd; }
  .content { padding: 20px; }
  .otp { font-size: 24px; font-weight: bold; color: #007bff; text-align: center; letter-spacing: 3px; margin: 20px 0; padding: 10px; background-color: #f2f2f2; border-radius: 3px; }
  .footer { font-size: 0.9em; text-align: center; color: #777; padding: 15px; border-top: 1px solid #ddd; }
</style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h2>{{.SenderName}} Email Verification</h2>
    </div>
    <div class="content">
      <p>Hello,</p>
      <p>Thank you for registering. Please use the following One-Time Password (OTP) to verify your email address:</p>
      <div class="otp">{{.OTP}}</div>
      <p>This code will expire in 15 minutes.</p>
      <p>If you did not request this, please ignore this email.</p>
    </div>
    <div class="footer">
      <p>&copy; {{.SenderName}}. All rights reserved.</p>
    </div>
  </div>
</body>
</html>
`

const passwordResetEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
<style>
  .container { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 20px auto; border: 1px solid #ddd; border-radius: 5px; }
  .header { background-color: #f7f7f7; padding: 15px; text-align: center; border-bottom: 1px solid #ddd; }
  .content { padding: 20px; }
  .otp { font-size: 24px; font-weight: bold; color: #dc3545; text-align: center; letter-spacing: 3px; margin: 20px 0; padding: 10px; background-color: #f2f2f2; border-radius: 3px; }
  .footer { font-size: 0.9em; text-align: center; color: #777; padding: 15px; border-top: 1px solid #ddd; }
  .warning { background-color: #fff3cd; padding: 10px; border-left: 4px solid #ffc107; margin: 15px 0; }
</style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h2>{{.SenderName}} Password Reset</h2>
    </div>
    <div class="content">
      <p>Hello,</p>
      <p>We received a request to reset your password. Please use the following One-Time Password (OTP) to proceed:</p>
      <div class="otp">{{.OTP}}</div>
      <p>This code will expire in 15 minutes.</p>
      <div class="warning">
        <strong>Security Notice:</strong> If you did not request a password reset, please ignore this email and ensure your account is secure.
      </div>
    </div>
    <div class="footer">
      <p>&copy; {{.SenderName}}. All rights reserved.</p>
    </div>
  </div>
</body>
</html>
`
