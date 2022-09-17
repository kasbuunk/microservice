package port

import "github.com/kasbuunk/microservice/app/email/models"

type Client interface {
	SendActivationLink(models.Address) error
}
