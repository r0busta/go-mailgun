package mail

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

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

// SendMessage sends a simple message with attachments
func (c *Client) SendMessage(from string, to string, subject string, text string, html string, attachments ...string) (string, error) {
	if text == "" && html == "" {
		return "", fmt.Errorf("empty text and html mail body")
	}

	m := c.mg.NewMessage(
		from,
		subject,
		text,
		to,
	)

	m.SetDKIM(true)

	if html != "" {
		m.SetHtml(html)
		m.SetTrackingOpens(true)
		m.SetTrackingClicks(true)
	}

	for _, f := range attachments {
		m.AddAttachment(f)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := c.mg.Send(ctx, m)
	return id, err
}
