package port

import "github.com/kasbuunk/microservice/app/email/models"

type EmailClient interface {
	SendActivationLink(models.Address) error
}
