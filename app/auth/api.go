package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt"

	"github.com/kasbuunk/microservice/app/auth/models"
	"github.com/kasbuunk/microservice/app/auth/user"
	"github.com/kasbuunk/microservice/app/client/eventbus"
	"github.com/kasbuunk/microservice/app/client/userrepo"
)

// API provides the interface that maps closely to however you wish to communicate with external components.
// It may be a one-to-one mapping to a graphql schema or grpc service.
// Other contexts, or 'domains', should communicate with each other through their APIs.
type API interface {
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

// Service implements the API. It has access to how its entities are stored and retrieved through its
// repositories. Additional repositories may be added here to access other entities. Other external clients are also
// added here so the domain core remains pure and agnostic of any calls over the network, including other
// microservices that are part of the same application.
type Service struct {
	UserRepo userrepo.Client
	EventBus eventbus.Client
}

func New(userRepo userrepo.Client, bus eventbus.Client) API {
	return Service{
		UserRepo: userRepo,
		EventBus: bus,
	}
}

func (s Service) Register(email models.EmailAddress, password models.Password) (models.User, error) {
	usr, err := user.NewUser(email, password)
	if err != nil {
		return usr, fmt.Errorf("saving user: %w", err)
	}

	savedUser, err := s.UserRepo.Save(usr)
	if err != nil {
		return savedUser, fmt.Errorf("saving user: %w", err)
	}

	// Invoke behaviour in Email service
	msg := eventbus.Event{
		Stream:  "AUTH",
		Subject: "USER_REGISTERED",
		Body:    eventbus.Body(fmt.Sprintf("new user registered with email %s", usr.Email)),
	}
	err = s.EventBus.Publish(msg)
	if err != nil {
		return savedUser, fmt.Errorf("publishing msg: %w", err)
	}
	return savedUser, nil
}

func (s Service) Login(email models.EmailAddress, password models.Password) (jwt.Token, error) {
	// TODO: Add method to repo to retrieve user by email, or add filter to List users and use the first row.
	// user := s.Users
	_, err := user.HashPassword(password)
	if err != nil {
		panic("handle me")
	}
	// if hash != user.PasswordHash {
	//	return jwt.Token{}, fmt.Errorf("password does not match")
	// }

	return jwt.Token{}, nil
}
