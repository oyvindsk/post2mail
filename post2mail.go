package post2mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strconv"
	"sync"
	"text/template"
)

// Templates for email the email sent
const defTemplate = `From: {{.FromEmail}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

	From:     {{.FromName}} ({{.FromEmail}})
    <br><br>
	Subject:  {{.Subject}}
    <br><br>
	Text:     {{.Text}}
    <br>
`

var onceInitTemplate sync.Once
var emailTemplate *template.Template // = template.Must(template.New("emailTemplate").Parse(defTemplateEmail))

// SMTPInfo contains the config params regarding the SMTP server
type SMTPInfo struct {
	Server   string
	Port     int
	Username string
	Password string
}

// EmailData contains the data to send in the email
type EmailData struct {
	FromName  string
	FromEmail string
	Subject   string
	Text      string
	To        string
}

// FormatAndSendEmail sends a simply formated email
func FormatAndSendEmail(ed EmailData, si SMTPInfo) error {

	// Parse the default email template once
	// caller can't change template ATM FIXME
	var err error
	onceInitTemplate.Do(func() { emailTemplate, err = template.New("").Parse(defTemplate) })
	if err != nil {
		return fmt.Errorf("FormatAndSendEmail: Parsing template failed: %s", err)
	}

	// Check that required fields are sat
	if si.Server == "" {
		return fmt.Errorf("FormatAndSendEmail: SMTPInfo.Server is required")
	}
	if ed.FromEmail == "" || ed.To == "" {
		return fmt.Errorf("FormatAndSendEmail: Some EmailData fields are required: To, FromEmail")
	}

	var buffer bytes.Buffer
	err = emailTemplate.Execute(&buffer, &ed)
	if err != nil {
		return fmt.Errorf("FormatAndSendEmail: could not execute email template: %s", err)
	}

	// log.Printf("Email:\n%q\n", buffer.String())

	auth := smtp.PlainAuth(
		"",
		si.Username,
		si.Password,
		si.Server,
	)

	return smtp.SendMail(
		si.Server+":"+strconv.Itoa(si.Port), // FIXME
		auth,
		si.Username,
		[]string{ed.To},
		buffer.Bytes(),
	)
}
