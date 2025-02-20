package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"

	"gopkg.in/gomail.v2"

	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
	"github.com/nathakusuma/elevateu-backend/internal/mailtmpl"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

type IMailer interface {
	Send(recipientEmail, subject, templateName string, data map[string]any) error
}

type mailer struct {
	dialer    *gomail.Dialer
	templates *template.Template
}

var (
	instance IMailer
	once     sync.Once
)

func NewMailDialer() IMailer {
	once.Do(func() {
		// Parse all templates at startup
		templates, err := template.ParseFS(mailtmpl.Templates, "*.html")
		if err != nil {
			log.Fatal(map[string]interface{}{
				"error": err.Error(),
			}, "[MAIL][NewMailDialer] failed to parse templates")
			return
		}

		instance = &mailer{
			dialer: gomail.NewDialer(
				env.GetEnv().SMTPHost,
				env.GetEnv().SMTPPort,
				env.GetEnv().SMTPUsername,
				env.GetEnv().SMTPPassword,
			),
			templates: templates,
		}
	})

	return instance
}

func (m *mailer) Send(recipientEmail, subject, templateName string, data map[string]any) error {
	var tmplOutput bytes.Buffer

	err := m.templates.ExecuteTemplate(&tmplOutput, templateName, data)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", "Vion <"+env.GetEnv().SMTPEmail+">")
	mail.SetHeader("To", recipientEmail)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", tmplOutput.String())

	return m.dialer.DialAndSend(mail)
}
