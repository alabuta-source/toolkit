package gomail

type emailOption struct {
	subject  string
	body     string
	replyTo  string
	file     string
	needCopy []string
	data     interface{}
}

type Option func(*emailOption)

func WithSubject(subject string) Option {
	return func(e *emailOption) {
		e.subject = subject
	}
}

func (e *emailOption) Subject() string {
	return e.subject
}

func WithBody(body string) Option {
	return func(e *emailOption) {
		e.body = body
	}
}

func (e *emailOption) Body() string {
	return e.body
}

func WithEmailToReplyTo(email string) Option {
	return func(e *emailOption) {
		e.replyTo = email
	}
}

func (e *emailOption) ReplyTo() string {
	return e.replyTo
}

func WithEmailCopyList(emails []string) Option {
	return func(e *emailOption) {
		e.needCopy = emails
	}
}

func (e *emailOption) NeedToCopy() []string {
	return e.needCopy
}

func WithData(data interface{}) Option {
	return func(e *emailOption) {
		e.data = data
	}
}

func (e *emailOption) Data() interface{} {
	return e.data
}

func WithHtmlFile(file string) Option {
	return func(e *emailOption) {
		e.file = file
	}
}

func (e *emailOption) File() string {
	return e.file
}
