package repository

import (
	"testing"

	"github.com/kasbuunk/microservice/api/auth"
	"github.com/kasbuunk/microservice/storage"
	"github.com/kasbuunk/microservice/test"
)

func getTestRepo(t *testing.T) auth.UserRepository {
	conf := storage.Config{
		Host: test.DBHost,
		Port: test.DBPort,
		Name: test.DBName,
		User: test.DBUser,
		Pass: test.DBPass,
	}
	db, err := storage.Connect(conf)
	if err != nil {
		t.Errorf("Connection to storage failed: %v", err)
	}
	repo := New(db)
	return repo
}

func cleanupRepo(t *testing.T, repo auth.UserRepository) {
	users, err := repo.List()
	if err != nil {
		t.Errorf("listing users: %v", err)
	}
	for _, user := range users {
		err := repo.Delete(user)
		if err != nil {
			t.Errorf("deleting user: %v", err)
		}
	}
}

func TestUserRepository_List(t *testing.T) {
	repo := getTestRepo(t)
	defer cleanupRepo(t, repo)

	// Check it does not have certain rows.
	_, err := repo.List()
	if err != nil {
		t.Errorf("listing users: %v", err)
	}
}

func TestUserRepository_Save(t *testing.T) {
	repo := getTestRepo(t)
	defer cleanupRepo(t, repo)

	user, err := auth.NewUser("Hex@amle.com", "jsdklfjsdlk2342A!")
	if err != nil {
		t.Errorf("instantiating new user: %v", err)
	}

	savedUser, err := repo.Save(user)
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

	newUser := auth.User{Email: "tnei@nrseit.sitm", PasswordHash: "rnsteinrisnter"}

	savedUser, err := repo.Save(newUser)
	if err != nil {
		t.Errorf("saving user: %v", err)
	}

	err = repo.Delete(savedUser)
	if err != nil {
		t.Errorf("deleting user: %v", err)
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo := getTestRepo(t)
	defer cleanupRepo(t, repo)

	newUser := auth.User{Email: "trstrs@nrseit.sitm", PasswordHash: "rstrstrs"}

	savedUser, err := repo.Save(newUser)
	if err != nil {
		t.Errorf("saving user: %v", err)
	}
	changedEmailField := auth.EmailAddress("mynew@email.address")
	savedUser.Email = changedEmailField
	changedUser, err := repo.Save(savedUser)
	if err != nil {
		t.Errorf("updating user: %v", err)
	}
	if changedUser.Email != changedEmailField {
		t.Errorf(test.ExpectedGot, changedEmailField, changedUser.Email)
	}
	if changedUser.ID != savedUser.ID {
		t.Error("ID changed on save")
	}
	if changedUser.PasswordHash != savedUser.PasswordHash {
		t.Error("Password changed on save")
	}

	err = repo.Delete(savedUser)
	if err != nil {
		t.Errorf("deleting user: %v", err)
	}
}

func TestUserRepository(t *testing.T) {
	repo := getTestRepo(t)
	defer cleanupRepo(t, repo)

	// Check it does not have certain rows.
	users, err := repo.List()
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
		user          auth.User
		expectedError error
	}{
		{
			"CorrectUser",
			auth.User{Email: "user@example.com", PasswordHash: "rmsetmremstrmsitmsreitm"},
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Insert rows.
			usr, err := repo.Save(tc.user)
			if err != tc.expectedError {
				t.Errorf(test.ExpectedGot, tc.expectedError, err)
			}
			// Check if returned user matches saved.
			if usr.Email != tc.user.Email {
				t.Errorf(test.ExpectedGot, tc.user.Email, usr.Email)
			}

			// Check rows exist after insertion.
			fetchedUser, err := repo.Load(usr.ID)
			if err != nil {
				t.Errorf("getting user from repo: %w", err)
			}
			if usr != fetchedUser {
				t.Errorf(test.ExpectedGot, usr, fetchedUser)
			}

			// Reset to original state.
			err = repo.Delete(usr)
			if err != nil {
				t.Errorf("deleting user: %v", err)
			}
			// Verify user does not exist anymore.
			_, err = repo.Load(usr.ID)
			if err == nil {
				t.Error("getting deleted user did not return error")
			}
		})
	}
}
