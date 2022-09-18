package core

import (
	"fmt"

	"github.com/golang-jwt/jwt"

	"github.com/kasbuunk/microservice/auth/models"
	"github.com/kasbuunk/microservice/auth/user"
	"github.com/kasbuunk/microservice/port"
)

// App implements the App. It has access to how its entities are stored and retrieved through its
// repositories. Additional repositories may be added here to access other entities. Other external clients are also
// added here so the domain core remains pure and agnostic of any calls over the network, including other
// microservices that are part of the same application.
type App struct {
	Repository port.Repository
	EventBus   port.EventBus
}

func New(userRepo port.Repository, bus port.EventBus) App {
	return App{
		Repository: userRepo,
		EventBus:   bus,
	}
}

func (a App) Register(email models.EmailAddress, password models.Password) (models.User, error) {
	usr, err := user.NewUser(email, password)
	if err != nil {
		return usr, fmt.Errorf("saving user: %w", err)
	}

	savedUser, err := a.Repository.UserSave(usr)
	if err != nil {
		return savedUser, fmt.Errorf("saving user: %w", err)
	}

	// Invoke behaviour in Email service
	msg := port.Event{
		Stream:  "AUTH",
		Subject: "USER_REGISTERED",
		Body:    port.Body(fmt.Sprintf("new user registered with email %s", usr.Email)),
	}
	err = a.EventBus.Publish(msg)
	if err != nil {
		return savedUser, fmt.Errorf("publishing msg: %w", err)
	}
	return savedUser, nil
}

func (a App) Login(email models.EmailAddress, password models.Password) (jwt.Token, error) {
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