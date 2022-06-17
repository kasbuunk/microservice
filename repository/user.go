package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/kasbuunk/microservice/auth"
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

func (us UserRepository) Save(usr auth.User) (auth.User, error) {
	// if id is set, update
	if usr.ID.String() != "" {
		updatedUser, err := updateByID(us.UserDB, tableName, usr)
		if err != nil {
			return usr, fmt.Errorf("updating user")
		}
		return updatedUser, nil
	}

	// if id is not set, create

	//if savedUser.ID.String() == "" {
	//	return &user.User{}, fmt.Errorf("saving user: %w", err)
	//}
	return usr, nil
}
func (us UserRepository) User(id uuid.UUID) (auth.User, error) {
	usr, err := selectByID(us.UserDB, tableName, id)
	if err != nil {
		return auth.User{}, fmt.Errorf("querying user by id: %w", err)
	}
	return usr, nil
}

func selectByID(db *sql.DB, table string, id uuid.UUID) (auth.User, error) {
	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT * FROM `%s` WHERE id = ? LIMIT 1;",
			table,
		),
		id.String(),
	)
	if err != nil {
		return auth.User{}, fmt.Errorf("querying db: %w", err)
	}

	// Candidate for generics. Replace the return 'user.User' type by a type parameter.
	var obj auth.User
	err = rows.Scan(obj)
	if err != nil {
		return auth.User{}, fmt.Errorf("scanning rows: %w", err)
	}

	return obj, nil
}

func updateByID(db *sql.DB, table string, obj auth.User) (auth.User, error) {
	result, err := db.ExecContext(context.Background(),
		fmt.Sprintf(
			"UPDATE `%s` SET `%s` WHERE id = ?;",
			table,
			updateSetValues([]string{
				"email",
			},
			),
		),
		string(obj.Email),
		obj.ID.String(),
	)
	if err != nil {
		return auth.User{}, fmt.Errorf("updating db: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return auth.User{}, fmt.Errorf(" rows: %w", err)
	}
	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", rows)
	}

	return obj, nil
}

func updateSetValues(fields []string) string {
	setString := "( "
	for _, field := range fields {
		setString = setString + fmt.Sprintf("%s = ?", field)
	}
	setString += " )"
	return setString
}
