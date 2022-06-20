package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

// Auth provides the api that maps closely to however you wish to communicate with external components.
// It may be a one-to-one mapping to a graphql schema or grpc service.
// Other contexts, or 'domains', should communicate with each other through their APIs.
type Auth interface {
	// Register inserts a new account into the repository, that has yet to be activated.
	Register(Email, Password) (User, error)
	Login(Email, Password) (jwt.Token, error)
	// ChangePassword(UserRepository, Password, Password, Password) (User, error)
	// ActivateAccount(UserRepository) (User, error)
	// Users takes a repo and a list of filters to return a list of Users.
	// Users(UserRepository, []string) ([]User, error)
	// User gets and returns a single User by ID from the repository.
	// User(id uuid.UUID) (User, error)
}

// Service implements the Auth interface. It has access to how Users are stored and retrieved through its
// UserRepository. Additional repositories may be added here to access other entities.
type Service struct {
	User UserRepository
}

func New(repo UserRepository) Auth {
	return Service{User: repo}
}

func (s Service) Register(email Email, password Password) (User, error) {
	user, err := NewUser(email, password)
	if err != nil {
		return user, fmt.Errorf("saving user: %w", err)
	}

	savedUser, err := s.User.Save(user)
	if err != nil {
		return savedUser, fmt.Errorf("saving user: %w", err)
	}
	return savedUser, nil
}

func (s Service) Login(email Email, password Password) (jwt.Token, error) {
	// TODO: Add method to repo to retrieve user by email, or add filter to List users and use the first row.
	// user := s.Users
	_, err := hashPassword(password)
	if err != nil {
		panic("handle me")
	}
	// if hash != user.PasswordHash {
	//	return jwt.Token{}, fmt.Errorf("password does not match")
	// }

	panic("Implement me!")
	return jwt.Token{}, nil
}
