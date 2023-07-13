package toolkit

import (
	"errors"
	"time"
)

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

func (payload *TokenPayload) valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New(formatErr(expiredTokenErr))
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

// GC_TYPE=
// GC_PROJECT_ID=
// GC_PRIVATE_KEY_ID=
// GC_PRIVATE_KEY=
// GC_CLIENT_EMAIL=
// GC_CLIENT_ID=
// GC_AUTH_URI=
// GC_TOKEN_URI=
// GC_AUTH_PROVIDER_X_CERT_URL=
// GC_CLIENT_X_CERT_URL=
// GC_UNIVERSE_DOMAIN=

type GCPBucketAuthJson struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}
