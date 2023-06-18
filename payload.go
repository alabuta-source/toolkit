package toolkit

import "time"

type TokenPayload struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(tokenID, username string, duration time.Duration) *TokenPayload {
	payload := &TokenPayload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload
}

func (payload *TokenPayload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return expiredTokenErr
	}
	return nil
}

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

// EmailTemplateBody is a struct that contains the configuration for the email template
// Use it to send data for simple templates that only need a name and an url.
type EmailTemplateBody struct {
	Name string
	URL  string
}
