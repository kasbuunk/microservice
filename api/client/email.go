package client

import "github.com/kasbuunk/microservice/api/email/models"

type EmailClient interface {
	SendActivationLink(models.Address) error
}
