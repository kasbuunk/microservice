package graph

import (
	"github.com/kasbuunk/microservice/app/auth"
)

// This file will not be regenerated automatically.
//
// It serves as eventbus injection for your app, add any dependencies you require here.

type Resolver struct {
	Auth auth.App
}
