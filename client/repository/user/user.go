package userrepo

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/kasbuunk/microservice/app/auth/models"
	"github.com/kasbuunk/microservice/app/dependency/userrepo"
)

const tableName = "users"

// UserRepository implements the Repository interface for the User model.
type UserRepository struct {
	UserDB *sql.DB
}

func New(db *sql.DB) userrepo.Client {
	return UserRepository{
		UserDB: db,
	}
}

func (us UserRepository) List() ([]models.User, error) {
	users, err := selectAll(us.UserDB, tableName)
	if err != nil {
		return []models.User{}, fmt.Errorf("listing users: %w", err)
	}
	return users, nil
}

func (us UserRepository) Delete(usr models.User) error {
	err := remove(us.UserDB, tableName, usr)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}

func (us UserRepository) Save(usr models.User) (models.User, error) {
	// If id is set, update, because the entity exists in storage.
	if !uuidIsEmpty(usr.ID) {
		err := update(us.UserDB, tableName, usr)
		if err != nil {
			return models.User{}, fmt.Errorf("updating user: %w", err)
		}
	} else {
		// Create (save an id-less user).
		err := insert(us.UserDB, tableName, usr)
		if err != nil {
			return models.User{}, fmt.Errorf("inserting user: %w", err)
		}
	}

	savedUser, err := selectByUniqueField(us.UserDB, tableName, "email", string(usr.Email))
	if err != nil {
		return models.User{}, fmt.Errorf("selecting inserted user: %w", err)
	}

	if uuidIsEmpty(savedUser.ID) {
		return models.User{}, fmt.Errorf("saved user has empty id")
	}
	return savedUser, nil
}

func (us UserRepository) Load(id uuid.UUID) (models.User, error) {
	usr, err := selectByID(us.UserDB, tableName, id)
	if err != nil {
		return models.User{}, fmt.Errorf("querying user by id: %w", err)
	}
	return usr, nil
}

func uuidIsEmpty(id uuid.UUID) bool {
	return id.ID() == 0
}
