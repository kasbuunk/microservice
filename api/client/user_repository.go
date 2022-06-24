package client

import (
	"github.com/google/uuid"

	"github.com/kasbuunk/microservice/api/auth/models"
)

// UserRepository is the storage interface through which User entities are retrieved, created and updated from storage.
type UserRepository interface {
	// Load 'gets' the instance from the repository by ID. Its implementation is abstracted away.
	Load(uuid.UUID) (models.User, error)
	// Save returns a new user instance. This is often achieved by passing in a pointer to a model instance,
	// with in-place modification of the object, like populating the ID generated by the database.
	// While that may be a little faster, returning a new instance has the benefits of immutable objects.
	// Any old reference will need the newly saved instance. Fewer side-effects, more expected behaviour.
	//
	// Create a new user by calling user.New(params), which checks any invariants for validation. Then, call its
	// Save method, which - if no ID is set yet - will save it to the chosen storage interface and populate its ID.
	//
	// Update an existing user by first getting it with the User(ID) method, calling methods that may change its
	// internal state and then call Save(user) to persist changes to the chosen storage interface.
	Save(models.User) (models.User, error)
	// List gets multiple users. TODO: accept a filter.
	List() ([]models.User, error)
	Delete(models.User) error
}