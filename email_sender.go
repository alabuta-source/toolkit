package toolkit

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
)

const (
	sendMailError = "Error sending email: [%s]"
	execTempErr   = "Error executing template: [%s]"
)

type EmailSender interface {
	// SendEmail sends a simple email.
	// to: the email address that will receive the email
	SendEmail(to string, subject string, body string) error

	// SendEmailWithSimpleTemplate sends an email using a simple a basic template.
	// templatePath: the path to the template file, the file must be a html file, example: "src/templates/email.html"
	// data: the data that will be used to fill the template, the template must have the following variables: {{.Name}} and {{.URL}}
	SendEmailWithSimpleTemplate(to string, subject string, templatePath string, data *EmailTemplateBody) error
}

type sender struct {
	Password     string
	From         string
	ServerConfig string
	Port         int
}

// NewEmailSender creates a new EmailSender instance.
// Instance it to access the EmailSender methods and send emails.
func NewEmailSender(configs *EmailSenderConfig) EmailSender {
	return &sender{
		Password:     configs.Password,
		From:         configs.From,
		ServerConfig: configs.ServerConfig,
		Port:         configs.Port,
	}
}

// SendEmail sends a simple email.
// to: the email address that will receive the email
func (s *sender) SendEmail(to string, subject string, body string) error {
	message := s.newMessage(to, subject, body)
	return s.dialAndSendMessage(message)
}

// SendEmailWithSimpleTemplate sends an email using a simple a basic template.
// templatePath: the path to the template file, the file must be a html file, example: "src/templates/email.html"
// data: the data that will be used to fill the template, the template must have the following variables: {{.Name}} and {{.URL}}
func (s *sender) SendEmailWithSimpleTemplate(to string, subject string, templatePath string, data *EmailTemplateBody) error {
	dir, err := getRootDir()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", dir, templatePath)
	temp, tErr := template.ParseFiles(path)
	if tErr != nil {
		return tErr
	}

	var bf bytes.Buffer
	if er := temp.Execute(&bf, data); er != nil {
		execErr := fmt.Sprintf(execTempErr, er.Error())
		return errors.New(execErr)
	}

	message := s.newMessage(to, subject, bf.String())
	return s.dialAndSendMessage(message)

}

func (s *sender) newMessage(to, subject, body string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", s.From)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	return message
}

func (s *sender) dialAndSendMessage(message *gomail.Message) error {
	dialer := gomail.NewDialer(s.ServerConfig, s.Port, s.From, s.Password)
	if err := dialer.DialAndSend(message); err != nil {
		dialerMessageErr := fmt.Sprintf(sendMailError, err.Error())
		return errors.New(dialerMessageErr)
	}
	return nil
}
