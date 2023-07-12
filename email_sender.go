package toolkit

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yosssi/gohtml"
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
	SendEmail(to string, subject string, body string, needCopy []string) error
	SendWelcomeEmail(to string, subject string, name string, needCopy []string) error
	SendResetPassEmail(to string, subject string, data *SimpleNotifyTemplate, needCopy []string) error
	SendVerifyEmail(to string, subject string, data *SimpleNotifyTemplate, needCopy []string) error
	SendBudgetEmail(to string, subject string, needCopy []string, data *SimpleNotifyTemplate) error
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
func (s *sender) SendEmail(to string, subject string, body string, needCopy []string) error {
	message := s.newMessage(to, subject, body, needCopy)
	return s.dialAndSendMessage(message)
}

func (s *sender) SendBudgetEmail(to string, subject string, needCopy []string, data *SimpleNotifyTemplate) error {
	return s.parseAndSend(to, subject, resetPassTemplate(), needCopy, data)
}

func (s *sender) SendWelcomeEmail(to string, subject string, name string, needCopy []string) error {
	return s.parseAndSend(to, subject, welcomeTemplate(), needCopy, &SimpleNotifyTemplate{Name: name})
}

func (s *sender) SendResetPassEmail(to string, subject string, data *SimpleNotifyTemplate, needCopy []string) error {
	return s.parseAndSend(to, subject, resetPassTemplate(), needCopy, data)
}

func (s *sender) SendVerifyEmail(to string, subject string, data *SimpleNotifyTemplate, needCopy []string) error {
	return s.parseAndSend(to, subject, verifyEmailTemplate(), needCopy, data)
}

func (s *sender) parseAndSend(to, subject string, file string, needCopy []string, data *SimpleNotifyTemplate) error {
	temp, tErr := template.New("toolkit_sender").Parse(file)
	if tErr != nil {
		return tErr
	}
	var bf bytes.Buffer
	if er := temp.Execute(gohtml.NewWriter(&bf), data); er != nil {
		execErr := fmt.Sprintf(execTempErr, er.Error())
		return errors.New(execErr)
	}
	message := s.newMessage(to, subject, bf.String(), needCopy)
	return s.dialAndSendMessage(message)
}

func (s *sender) newMessage(to string, subject, body string, needCopy []string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", s.From)
	message.SetHeader("To", to)
	message.SetHeader("Cc", needCopy...)
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
