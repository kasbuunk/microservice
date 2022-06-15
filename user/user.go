package user

import (
	"fmt"

	"github.com/google/uuid"
)

type Email string
type PasswordHash string
type Password string

type User struct {
	ID           uuid.UUID
	Email        Email
	PasswordHash PasswordHash
}

// New validates a user's invariants and returns the User. It does not save the entity to storage.
// The caller must persist the new entity by calling its `Save()` method.
func New(email, password string) (User, error) {
	var user User

	emailX := Email(email)
	passwX := Password(email)

	err := validateEmail(emailX)
	if err != nil {
		return user, err
	}
	passwordHash, err := validatePassword(passwX)
	if err != nil {
		return user, err
	}

	user.Email = emailX
	user.PasswordHash = passwordHash

	return user, nil
}

func validateEmail(email Email) error {
	return fmt.Errorf("email '%s' invalid", email)
}

func validatePassword(_ Password) (PasswordHash, error) {
	return "", fmt.Errorf("password invalid")
}
