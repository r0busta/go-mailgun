package mail

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/r0busta/mailgun-go/v4"
)

type Message struct {
	From    string
	To      string
	Subject string
	Text    string
	HTML    string

	Attachments       []string
	BufferAttachments []BufferAttachment

	Inlines       []string
	BufferInlines []BufferAttachment
}

type BufferAttachment struct {
	Filename string
	Buffer   []byte
}

// Client mail client
type Client struct {
	mg *mailgun.MailgunImpl
}

// NewDefaultClient creates new mail client
func NewDefaultClient() (*Client, error) {
	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")

	if domain == "" || apiKey == "" {
		return nil, fmt.Errorf("Mailgun API key and domain not configured")
	}

	c := &Client{}
	mgClient := mailgun.NewMailgun(domain, apiKey)
	mgClient.SetAPIBase(mailgun.APIBaseEU)

	c.mg = mgClient

	return c, nil
}

// SendMessage a convenient function to send a simple message with attachments
func (c *Client) SendMessage(m *Message) (string, error) {
	if m.Text == "" && m.HTML == "" {
		return "", fmt.Errorf("empty text and html mail body")
	}

	mgm := c.mg.NewMessage(
		m.From,
		m.Subject,
		m.Text,
		m.To,
	)

	mgm.SetDKIM(true)

	if m.HTML != "" {
		mgm.SetHtml(m.HTML)
		mgm.SetTrackingOpens(true)
		mgm.SetTrackingClicks(true)
	}

	for _, a := range m.Attachments {
		mgm.AddAttachment(a)
	}

	for _, a := range m.Inlines {
		mgm.AddInline(a)
	}

	for _, a := range m.BufferAttachments {
		mgm.AddBufferAttachment(a.Filename, a.Buffer)
	}

	for _, a := range m.BufferInlines {
		mgm.AddBufferInline(a.Filename, a.Buffer)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := c.mg.Send(ctx, mgm)
	return id, err
}
