package auth

import (
	"github.com/golang-jwt/jwt"

	"github.com/kasbuunk/microservice/app/auth/models"
)

// App provides the interface that maps closely to however you wish to communicate with external components.
// It may be a one-to-one mapping to a graphql schema or grpc service.
// Other contexts, or 'domains', should communicate with each other through their APIs.
type App interface {
	// Register inserts a new account into the repository, that has yet to be activated.
	Register(models.EmailAddress, models.Password) (models.User, error)
	Login(models.EmailAddress, models.Password) (jwt.Token, error)
	// ChangePassword(UserRepository, Password, Password, Password) (User, error)
	// ActivateAccount(UserRepository) (User, error)
	// Users takes a repo and a list of filters to return a list of Users.
	// Users(UserRepository, []string) ([]User, error)
	// User gets and returns a single User by ID from the repository.
	// User(id uuid.UUID) (User, error)
}
