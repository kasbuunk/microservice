package userrepo

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/kasbuunk/microservice/auth/models"
)

const tableName = "users"

// Repository implements the UserRepository interface for the User model.
type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}

func (us Repository) Users() ([]models.User, error) {
	users, err := selectAll(us.DB, tableName)
	if err != nil {
		return []models.User{}, fmt.Errorf("listing users: %w", err)
	}
	return users, nil
}

func (us Repository) UserDelete(usr models.User) error {
	err := remove(us.DB, tableName, usr)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}

func (us Repository) UserSave(usr models.User) (models.User, error) {
	// If id is set, update, because the entity exists in storage.
	if !uuidIsEmpty(usr.ID) {
		err := update(us.DB, tableName, usr)
		if err != nil {
			return models.User{}, fmt.Errorf("updating user: %w", err)
		}
	} else {
		// Create (save an id-less user).
		err := insert(us.DB, tableName, usr)
		if err != nil {
			return models.User{}, fmt.Errorf("inserting user: %w", err)
		}
	}

	savedUser, err := selectByUniqueField(us.DB, tableName, "email", string(usr.Email))
	if err != nil {
		return models.User{}, fmt.Errorf("selecting inserted user: %w", err)
	}

	if uuidIsEmpty(savedUser.ID) {
		return models.User{}, fmt.Errorf("saved user has empty id")
	}
	return savedUser, nil
}

func (us Repository) User(id uuid.UUID) (models.User, error) {
	usr, err := selectByID(us.DB, tableName, id)
	if err != nil {
		return models.User{}, fmt.Errorf("querying user by id: %w", err)
	}
	return usr, nil
}

func uuidIsEmpty(id uuid.UUID) bool {
	return id.ID() == 0
}
