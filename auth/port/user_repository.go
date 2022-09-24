package port

import (
	"github.com/google/uuid"

	"github.com/kasbuunk/microservice/auth/models"
)

// Repository is the interface through which User entities are retrieved, created and updated from storage.
type Repository interface {
	// User 'gets' the instance from the repository by ID. Its implementation is abstracted away.
	User(uuid.UUID) (models.User, error)
	// UserSave returns a new user instance. This is often achieved by passing in a pointer to a model instance,
	// with in-place modification of the object, like populating the ID generated by the database.
	// While that may be a little faster, returning a new instance has the benefits of immutable objects.
	// Any old reference will need the newly saved instance. Fewer side effects, more expected behaviour.
	//
	// Create a new user by calling user.New(params), which checks any invariants for validation. Then, call its
	// Save method, which - if no ID is set yet - will save it to the chosen storage interface and populate its ID.
	//
	// Save an existing user by first getting it with the User(ID) method, calling methods that may change its
	// internal state and then call Save(user) to persist changes to the chosen storage interface.
	UserSave(models.User) (models.User, error)
	// Users gets multiple users. TODO: accept a filter.
	Users() ([]models.User, error)
	UserDelete(models.User) error
}