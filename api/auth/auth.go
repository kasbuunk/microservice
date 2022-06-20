package auth

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"

	"github.com/kasbuunk/microservice/events"
)

// Auth provides the api that maps closely to however you wish to communicate with external components.
// It may be a one-to-one mapping to a graphql schema or grpc service.
// Other contexts, or 'domains', should communicate with each other through their APIs.
type Auth interface {
	Subscribe()
	// Register inserts a new account into the repository, that has yet to be activated.
	Register(EmailAddress, Password) (User, error)
	Login(EmailAddress, Password) (jwt.Token, error)
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
	User UserRepository
	Bus  events.MessageBus
}

func New(repo UserRepository, bus events.MessageBus) Auth {
	return Service{
		User: repo,
		Bus:  bus,
	}
}

// Subscribe listens for messages that match the Stream or Subject.
func (s Service) Subscribe() {
	fmt.Println("listening for messages")
	// starts process in loop, in goroutine that awaits published messages and invokes api calls
	messageChannel, err := s.Bus.Subscribe("EMAIL", "ACTIVATION_REQUEST_SENT")
	if err != nil {
		log.Fatal(fmt.Errorf("subscribing: %w", err))
	}
	fmt.Println(<-messageChannel)

}

func (s Service) Register(email EmailAddress, password Password) (User, error) {
	user, err := NewUser(email, password)
	if err != nil {
		return user, fmt.Errorf("saving user: %w", err)
	}

	savedUser, err := s.User.Save(user)
	if err != nil {
		return savedUser, fmt.Errorf("saving user: %w", err)
	}

	// Invoke behaviour in Email service
	msg := events.Message{
		Stream:  "AUTH",
		Subject: "USER_REGISTERED",
		Body:    fmt.Sprintf("new user registered with email %s", user.Email),
	}
	err = s.Bus.Publish(msg)
	if err != nil {
		return savedUser, fmt.Errorf("publishing msg: %w", err)
	}
	return savedUser, nil
}

func (s Service) Login(email EmailAddress, password Password) (jwt.Token, error) {
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