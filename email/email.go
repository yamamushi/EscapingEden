package email

import "github.com/yamamushi/EscapingEden/messages"

type EmailProviderID int // Enum

const (
	EmailProvider_SMTP EmailProviderID = iota
	EmailProvider_Gmail
)

func (eid *EmailProviderID) String() string {
	switch *eid {
	case EmailProvider_SMTP:
		return "SMTP"
	case EmailProvider_Gmail:
		return "Gmail"
	default:
		return "Unknown"
	}
}

type EmailProviderType interface {
	GetID() EmailProviderID

	SendEmail(email messages.EmailMessage) error
}

type EmailProvider struct {
	ID EmailProviderID
}

func NewEmailProvider(id EmailProviderID) *EmailProvider {
	return &EmailProvider{ID: id}
}

func (p *EmailProvider) GetID() EmailProviderID {
	return p.ID
}

func (p *EmailProvider) SendEmail(email messages.EmailMessage) error {
	return nil
}
