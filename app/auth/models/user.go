package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        EmailAddress
	PasswordHash PasswordHash
}

type EmailAddress string
type PasswordHash string
type Password string
