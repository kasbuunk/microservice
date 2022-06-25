package email

import "github.com/kasbuunk/microservice/api/email/models"

type Client interface {
	SendActivationLink(models.Address) error
}
