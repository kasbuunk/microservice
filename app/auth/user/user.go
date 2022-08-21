package user

import (
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"

	"github.com/kasbuunk/microservice/app/auth/models"
)

const EmailMaxLength = 64

// NewUser validates a user's invariants and returns the User. It does not save the entity to storage.
// The caller must persist the new entity by calling its `Save()` method.
func NewUser(email models.EmailAddress, password models.Password) (models.User, error) {
	var user models.User

	err := ValidateEmail(email)
	if err != nil {
		return user, err
	}
	err = ValidatePassword(password)
	if err != nil {
		return user, err
	}
	passwordHash, err := HashPassword(password)
	if err != nil {
		return user, err
	}

	user.Email = email
	user.PasswordHash = passwordHash

	return user, nil
}

func ValidateEmail(email models.EmailAddress) error {
	_, err := mail.ParseAddress(string(email))
	if err != nil {
		return fmt.Errorf("email '%s' invalid", email)
	}
	if len(email) > EmailMaxLength {
		return fmt.Errorf("email cannot have more than 64 characters")
	}
	return nil
}

func ValidatePassword(_ models.Password) error {
	return nil
}

func HashPassword(password models.Password) (models.PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", fmt.Errorf("password invalid")
	}
	return models.PasswordHash(hash), nil
}
