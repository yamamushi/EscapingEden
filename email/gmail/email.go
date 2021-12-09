package gmail

import (
	"github.com/yamamushi/EscapingEden/email"
	"github.com/yamamushi/EscapingEden/messages"
)

type GmailProvider struct {
	email.EmailProviderType

	Username string
	Password string
}

func NewGmailProvider(username, password string) (*GmailProvider, error) {
	return &GmailProvider{
		Username: username,
		Password: password,
	}, nil
}

func (g *GmailProvider) SendEmail(email messages.EmailMessage) error {
	return nil
}
