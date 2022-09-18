package userrepo

import (
	"github.com/kasbuunk/microservice/app/adapter/repository/storage"
	"github.com/kasbuunk/microservice/app/port"
	"testing"

	"github.com/kasbuunk/microservice/app/auth/models"
	"github.com/kasbuunk/microservice/app/auth/user"
)

var (
	DBHost = "localhost"
	DBPort = 5432
	DBName = "auth_test"
	DBUser = "postgres"
	DBPass = "postgres"

	ExpectedGot = "expected '%v'; got '%v'"
)

func getTestRepo(t *testing.T) port.Repository {
	conf := storage.Config{
		Host: DBHost,
		Port: DBPort,
		Name: DBName,
		User: DBUser,
		Pass: DBPass,
	}
	db, err := storage.Connect(conf)
	if err != nil {
		t.Errorf("Connection to storage failed: %v", err)
	}
	repo := New(db)
	return repo
}

func cleanupRepo(t *testing.T, repo port.Repository) {
	users, err := repo.Users()
	if err != nil {
		t.Errorf("listing users: %v", err)
	}
	for _, usr := range users {
		err := repo.UserDelete(usr)
		if err != nil {
			t.Errorf("deleting user: %v", err)
		}
	}
}

func TestUserRepository_List(t *testing.T) {
	repo := getTestRepo(t)
	defer cleanupRepo(t, repo)

	// Check it does not have certain rows.
	_, err := repo.Users()
	if err != nil {
		t.Errorf("listing users: %v", err)
	}
}

func TestUserRepository_Save(t *testing.T) {
	repo := getTestRepo(t)
	defer cleanupRepo(t, repo)

	usr, err := user.NewUser("Hex@amle.com", "jsdklfjsdlk2342A!")
	if err != nil {
		t.Errorf("instantiating new user: %v", err)
	}

	savedUser, err := repo.UserSave(usr)
	if err != nil {
		t.Errorf("saving user: %v", err)
	}
	if uuidIsEmpty(savedUser.ID) {
		t.Errorf("empty id after save, postgres should auto generate uuids")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo := getTestRepo(t)
	defer cleanupRepo(t, repo)

	newUser := models.User{Email: "tnei@nrseit.sitm", PasswordHash: "rnsteinrisnter"}

	savedUser, err := repo.UserSave(newUser)
	if err != nil {
		t.Errorf("saving user: %v", err)
	}

	err = repo.UserDelete(savedUser)
	if err != nil {
		t.Errorf("deleting user: %v", err)
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo := getTestRepo(t)
	defer cleanupRepo(t, repo)

	newUser := models.User{Email: "trstrs@nrseit.sitm", PasswordHash: "rstrstrs"}

	savedUser, err := repo.UserSave(newUser)
	if err != nil {
		t.Errorf("saving user: %v", err)
	}
	changedEmailField := models.EmailAddress("mynew@email.address")
	savedUser.Email = changedEmailField
	changedUser, err := repo.UserSave(savedUser)
	if err != nil {
		t.Errorf("updating user: %v", err)
	}
	if changedUser.Email != changedEmailField {
		t.Errorf(ExpectedGot, changedEmailField, changedUser.Email)
	}
	if changedUser.ID != savedUser.ID {
		t.Error("ID changed on save")
	}
	if changedUser.PasswordHash != savedUser.PasswordHash {
		t.Error("Password changed on save")
	}

	err = repo.UserDelete(savedUser)
	if err != nil {
		t.Errorf("deleting user: %v", err)
	}
}

func TestUserRepository(t *testing.T) {
	repo := getTestRepo(t)
	defer cleanupRepo(t, repo)

	// Check it does not have certain rows.
	users, err := repo.Users()
	if err != nil {
		t.Errorf("listing users: %v", err)
	}

	if len(users) > 0 {
		// When tests run in parallel and use the same database, this might be flaky, since other
		// tests produce side-effects in the same state.
		t.Errorf("%d users found, but should be 0", len(users))
	}

	cases := []struct {
		name          string
		user          models.User
		expectedError error
	}{
		{
			"CorrectUser",
			models.User{Email: "user@example.com", PasswordHash: "rmsetmremstrmsitmsreitm"},
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Insert rows.
			usr, err := repo.UserSave(tc.user)
			if err != tc.expectedError {
				t.Errorf(ExpectedGot, tc.expectedError, err)
			}
			// Check if returned user matches saved.
			if usr.Email != tc.user.Email {
				t.Errorf(ExpectedGot, tc.user.Email, usr.Email)
			}

			// Check rows exist after insertion.
			fetchedUser, err := repo.User(usr.ID)
			if err != nil {
				t.Errorf("getting user from repo: %v", err)
			}
			if usr != fetchedUser {
				t.Errorf(ExpectedGot, usr, fetchedUser)
			}

			// Reset to original state.
			err = repo.UserDelete(usr)
			if err != nil {
				t.Errorf("deleting user: %v", err)
			}
			// Verify user does not exist anymore.
			_, err = repo.User(usr.ID)
			if err == nil {
				t.Error("getting deleted user did not return error")
			}
		})
	}
}
