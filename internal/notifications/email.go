package notifications

import (
	"fmt"
	"net"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SimpleEmail struct {
	To      string
	Subject string
	Body    string
}

type EmailNotifier struct {
	config *SMTPConfig
}

func NewEmailNotifier(config *SMTPConfig) *EmailNotifier {
	return &EmailNotifier{config: config}
}

func (e *EmailNotifier) SendSimpleEmail(email *SimpleEmail) error {
	addr := net.JoinHostPort(e.config.Host, fmt.Sprintf("%d", e.config.Port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer func() {
		if err = conn.Close(); err != nil {
			fmt.Println("Failed to close connection:", err)
		}
	}()

	client, smtpErr := smtp.NewClient(conn, e.config.Host)
	if smtpErr != nil {
		return smtpErr
	}
	defer func() {
		if err = client.Close(); err != nil {
			fmt.Println("Failed to close SMTP client:", err)
		}
	}()

	if e.config.Username != "" || e.config.Password != "" {
		if authErr := client.Auth(smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)); authErr != nil {
			return authErr
		}
	}

	if mailErr := client.Mail(e.config.From); mailErr != nil {
		return mailErr
	}

	if rcptErr := client.Rcpt(email.To); rcptErr != nil {
		return rcptErr
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("To: %s\r\n", email.To)
	msg += fmt.Sprintf("Subject: %s\r\n", email.Subject)
	msg += fmt.Sprintf("From: %s\r\n", e.config.From)
	msg += fmt.Sprintf("\r\n%s\r\n", email.Body)

	if _, err := writer.Write([]byte(msg)); err != nil {
		return err
	}

	return writer.Close()
}

func (e *EmailNotifier) SendLoginNotification(email, name string) error {
	return e.SendSimpleEmail(&SimpleEmail{
		To:      email,
		Subject: "Login Notification",
		Body:    fmt.Sprintf("Hello %s, you have successfully logged in to your account. If you did not make this request, please contact support.", name),
	})
}
