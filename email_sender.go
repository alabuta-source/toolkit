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

	SendBudgetEmail(to string, subject string, file string, data *SimpleNotifyTemplate) error
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

func (s *sender) SendBudgetEmail(to string, subject string, file string, data *SimpleNotifyTemplate) error {
	return s.parseAndSend(to, subject, file, data)
}

func (s *sender) parseAndSend(to, subject, file string, data *SimpleNotifyTemplate) error {
	path := fmt.Sprintf("%s/%s", "templates", file)
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
