package user

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/kasbuunk/microservice/api/auth"
)

const tableName = "users"

// UserRepository implements the Repository interface for the User model.
type UserRepository struct {
	UserDB *sql.DB
}

func New(db *sql.DB) auth.UserRepository {
	return UserRepository{
		UserDB: db,
	}
}

func (us UserRepository) List() ([]auth.User, error) {
	users, err := selectAll(us.UserDB, tableName)
	if err != nil {
		return []auth.User{}, fmt.Errorf("listing users: %w", err)
	}
	return users, nil
}

func (us UserRepository) Delete(usr auth.User) error {
	err := remove(us.UserDB, tableName, usr)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}

func (us UserRepository) Save(usr auth.User) (auth.User, error) {
	// If id is set, update, because the entity exists in storage.
	if !uuidIsEmpty(usr.ID) {
		err := update(us.UserDB, tableName, usr)
		if err != nil {
			return auth.User{}, fmt.Errorf("updating user: %w", err)
		}
	} else {
		// Create (save an id-less user).
		err := insert(us.UserDB, tableName, usr)
		if err != nil {
			return auth.User{}, fmt.Errorf("inserting user: %w", err)
		}
	}

	savedUser, err := selectByUniqueField(us.UserDB, tableName, "email", string(usr.Email))
	if err != nil {
		return auth.User{}, fmt.Errorf("selecting inserted user: %w", err)
	}

	if uuidIsEmpty(savedUser.ID) {
		return auth.User{}, fmt.Errorf("saved user has empty id")
	}
	return savedUser, nil
}

func (us UserRepository) Load(id uuid.UUID) (auth.User, error) {
	usr, err := selectByID(us.UserDB, tableName, id)
	if err != nil {
		return auth.User{}, fmt.Errorf("querying user by id: %w", err)
	}
	return usr, nil
}

func uuidIsEmpty(id uuid.UUID) bool {
	if id.ID() == 0 {
		return true
	}
	return false
}
