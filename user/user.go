package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/mail"

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
func New(email Email, password Password) (User, error) {
	var user User

	err := validateEmail(email)
	if err != nil {
		return user, err
	}
	err = validatePassword(password)
	passwordHash, err := hashPassword(password)
	if err != nil {
		return user, err
	}

	user.Email = email
	user.PasswordHash = passwordHash

	return user, nil
}

func validateEmail(email Email) error {
	_, err := mail.ParseAddress(string(email))
	if err != nil {
		return fmt.Errorf("email '%s' invalid", email)
	}
	if len(email) > 64 {
		return fmt.Errorf("email cannot have more than 64 characters")
	}
	return nil
}

func validatePassword(_ Password) error {
	return fmt.Errorf("password invalid")
}

func hashPassword(password Password) (PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", fmt.Errorf("password invalid")
	}
	return PasswordHash(hash), nil
}
