package core

import (
	"fmt"

	"github.com/golang-jwt/jwt"

	"github.com/kasbuunk/microservice/app/auth"
	"github.com/kasbuunk/microservice/app/auth/dependency"
	"github.com/kasbuunk/microservice/app/auth/models"
	"github.com/kasbuunk/microservice/app/auth/user"
	"github.com/kasbuunk/microservice/app/eventbus"
)

// application implements the App. It has access to how its entities are stored and retrieved through its
// repositories. Additional repositories may be added here to access other entities. Other external clients are also
// added here so the domain core remains pure and agnostic of any calls over the network, including other
// microservices that are part of the same application.
type application struct {
	UserRepo dependency.UserRepository
	EventBus eventbus.Client
}

func New(userRepo dependency.UserRepository, bus eventbus.Client) auth.App {
	return application{
		UserRepo: userRepo,
		EventBus: bus,
	}
}

func (a application) Register(email models.EmailAddress, password models.Password) (models.User, error) {
	usr, err := user.NewUser(email, password)
	if err != nil {
		return usr, fmt.Errorf("saving user: %w", err)
	}

	savedUser, err := a.UserRepo.Save(usr)
	if err != nil {
		return savedUser, fmt.Errorf("saving user: %w", err)
	}

	// Invoke behaviour in Email service
	msg := eventbus.Event{
		Stream:  "AUTH",
		Subject: "USER_REGISTERED",
		Body:    eventbus.Body(fmt.Sprintf("new user registered with email %s", usr.Email)),
	}
	err = a.EventBus.Publish(msg)
	if err != nil {
		return savedUser, fmt.Errorf("publishing msg: %w", err)
	}
	return savedUser, nil
}

func (a application) Login(email models.EmailAddress, password models.Password) (jwt.Token, error) {
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
