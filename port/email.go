package port

import (
	"github.com/kasbuunk/microservice/email/models"
)

type EmailClient interface {
	SendActivationLink(models.Address) error
}
