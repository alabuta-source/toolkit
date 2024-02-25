package gomail

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yosssi/gohtml"
	"gopkg.in/gomail.v2"
	"html/template"
)

// EmailSenderConfig is a struct that contains the configuration for the email sender
// Use it to create a new email sender client passing it as a parameter
type EmailSenderConfig struct {
	// The email account secret password
	Password string
	// The email account owner identification, this is the email address that will be used to send the emails
	From string
	// The email server configuration, example: smtp.gmail.com
	ServerConfig string
	// The email server port
	Port int
}

// EmailTemplateBody
// Use it to send data for simple templates that only need a name and an url.
type EmailTemplateBody struct {
	Name string
	URL  string
}

type SimpleNotifyTemplate struct {
	Name    string
	Message string
}

type BudgetTemplateBody struct {
	Name    string
	Message string
	Phone   string
	Email   string
}

const (
	sendMailError = "Error sending email: [%s]"
	execTempErr   = "Error executing template: [%s]"
)

type EmailSender interface {
	// SendEmail sends a simple email.
	// to: the email address that will receive the email
	SendEmail(to string, option ...Option) error
	SendWelcomeEmail(to string, option ...Option) error
	SendResetPassEmail(to string, option ...Option) error
	SendVerifyEmail(to string, option ...Option) error
	SendBudgetEmail(to string, option ...Option) error
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
func (s *sender) SendEmail(to string, option ...Option) error {
	reqOptions := s.setupOptions(option...)
	message := s.newMessage(to, reqOptions.Subject(), reqOptions.Body(), reqOptions.NeedToCopy())
	return s.dialAndSendMessage(message)
}

func (s *sender) SendBudgetEmail(to string, option ...Option) error {
	reqOptions := s.setupOptions(option...)
	return s.parseAndSend(to, reqOptions.Subject(), budgetTemplate(), reqOptions.NeedToCopy(), reqOptions.Data())
}

func (s *sender) SendWelcomeEmail(to string, option ...Option) error {
	reqOptions := s.setupOptions(option...)
	return s.parseAndSend(to, reqOptions.Subject(), welcomeTemplate(), reqOptions.NeedToCopy(), reqOptions.Data())
}

func (s *sender) SendResetPassEmail(to string, option ...Option) error {
	reqOptions := s.setupOptions(option...)
	return s.parseAndSend(to, reqOptions.Subject(), resetPassTemplate(), reqOptions.NeedToCopy(), reqOptions.Data())
}

func (s *sender) SendVerifyEmail(to string, option ...Option) error {
	reqOptions := s.setupOptions(option...)
	return s.parseAndSend(to, reqOptions.Subject(), verifyEmailTemplate(), reqOptions.NeedToCopy(), reqOptions.Data())
}

func (s *sender) parseAndSend(to, subject string, file string, needCopy []string, data interface{}) error {
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

func (*sender) setupOptions(option ...Option) *emailOption {
	var reqOptions emailOption
	for _, opt := range option {
		opt(&reqOptions)
	}
	return &reqOptions
}
