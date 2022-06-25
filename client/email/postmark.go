// Package emailclient implements the EmailClient interface, so the core domain can remain agnostic of
// the implementation of how emails are actually sent when all domain logic is established.
package emailclient

import (
	"fmt"

	"github.com/keighl/postmark"

	"github.com/kasbuunk/microservice/api/client/email"
	"github.com/kasbuunk/microservice/api/email/models"
)

type Config struct {
	ServerToken  string
	AccountToken string
}

// New returns a configured email client.
// TODO: Add static configuration in config and pass through here as input.
func New(conf Config) email.Client {
	postmarkClient := postmark.NewClient(conf.ServerToken, conf.AccountToken)

	return EmailClient{
		Client: *postmarkClient,
	}
}

type EmailClient struct {
	Client postmark.Client
}

func (ec EmailClient) SendActivationLink(userEmailAddress models.Address) error {
	msg := postmark.Email{
		From:        "no-reply@example.com",
		To:          string(userEmailAddress),
		Subject:     "Activate your account",
		HtmlBody:    "...",
		TextBody:    "...",
		Tag:         "activate-account",
		TrackOpens:  true,
		Cc:          "",
		Bcc:         "",
		ReplyTo:     "",
		Headers:     nil,
		Attachments: nil,
		Metadata:    nil,
	}
	response, err := ec.Client.SendEmail(msg)
	if err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	// TODO: Do things with response.
	_ = response
	return nil
}
