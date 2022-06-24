package userrepo

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"

	"github.com/kasbuunk/microservice/api/auth/models"
)

const insertPrefix = "INSERT INTO"
const updatePrefix = "UPDATE"

const insertFmt = "(%s) VALUES (%s)"
const updateFmt = "SET %s WHERE %s"

type fieldValues map[string]string

func insertQuery(table string, fields fieldValues) string {
	var columns string
	var values string

	for field, value := range fields {
		columns += field + ", "
		values += value + ", "
	}

	trimmedColumns := strings.Trim(columns, " ,")
	trimmedValues := strings.Trim(values, " ,")

	queryBody := fmt.Sprintf(insertFmt, trimmedColumns, trimmedValues)
	queryString := fmt.Sprintf("%s %s %s;", insertPrefix, table, queryBody)

	return queryString
}

func updateQuery(table string, fields fieldValues, where string) string {
	var setString string
	for field, value := range fields {
		setString += field + " = "
		setString += value + ", "
	}

	trimmedString := strings.Trim(setString, " ,")

	queryBody := fmt.Sprintf(updateFmt, trimmedString, where)
	queryString := fmt.Sprintf("%s %s %s;", updatePrefix, table, queryBody)

	return queryString
}

func insert(db *sql.DB, table string, usr models.User) error {
	fieldValues := map[string]string{
		"email":         "$1",
		"password_hash": "$2",
	}
	query := insertQuery(table, fieldValues)
	err := executeQuery(db, query, string(usr.Email), string(usr.PasswordHash))
	if err != nil {
		return fmt.Errorf("inserting: %w", err)
	}
	return nil
}

func remove(db *sql.DB, table string, usr models.User) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", table)
	err := executeQuery(db, query, usr.ID)
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}
	return nil
}

func update(db *sql.DB, table string, usr models.User) error {
	fieldValues := map[string]string{
		"email":         "$1",
		"password_hash": "$2",
	}
	query := updateQuery(table, fieldValues, "id = $3")
	err := executeQuery(db, query, string(usr.Email), string(usr.PasswordHash), usr.ID.String())
	if err != nil {
		return fmt.Errorf("updating: %w", err)
	}
	return nil
}

func executeQuery(db *sql.DB, query string, params ...interface{}) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("Error: closing statement: %v", err)
		}
	}(stmt)

	res, err := stmt.Exec(params...)
	if err != nil {
		return fmt.Errorf("executing statement: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}
	return nil
}

func selectAll(db *sql.DB, table string) ([]models.User, error) {
	var users []models.User

	query := fmt.Sprintf(
		"SELECT id, email, password_hash FROM %s;",
		table,
	)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("selecting all: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error: closing rows: %v", err)
			return
		}
	}(rows)

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash)
		if err != nil {
			return nil, fmt.Errorf("scanning rows: %w", err)
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("iterating over results: %w", err)
	}

	return users, nil
}

func selectByUniqueField(db *sql.DB, table, field, value string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s = $1;",
		table,
		field,
	)
	// Query for a value based on a single row.
	err := db.QueryRow(query, value).Scan(&user.ID, &user.Email, &user.PasswordHash)
	switch err {
	case nil:
		return user, nil
	case sql.ErrNoRows:
		return user, err
	default:
		return user, fmt.Errorf("selecting user: %w", err)
	}
}

func selectByID(db *sql.DB, table string, id uuid.UUID) (models.User, error) {
	user, err := selectByUniqueField(db, table, "id", id.String())
	if err != nil {
		return models.User{}, fmt.Errorf("selecting id: %w", err)
	}
	return user, nil
}
