package storage

import (
	"errors"
	"testing"
)

var (
	DBHost = "localhost"
	DBPort = 5432
	DBName = "auth_test"
	DBUser = "postgres"
	DBPass = "postgres"
)

func getTestDB() Config {
	return Config{
		Host: DBHost,
		Port: DBPort,
		Name: DBName,
		User: DBUser,
		Pass: DBPass,
	}
}

// Given: test database is running on localhost with the above configuration.
func TestConnectDatabase(t *testing.T) {
	// Connect to database. Includes a ping, so no further verification needed.
	_, err := Connect(getTestDB())
	if err != nil {
		t.Log(getTestDB())
		t.Error(err)
	}
}

// Given: test database is running on localhost:5432 with user/pass postgres.
func TestConnectNonExistingDatabase(t *testing.T) {
	dbFaultyConf := getTestDB()
	dbFaultyConf.Name = "faulty_override"

	_, err := Connect(dbFaultyConf)
	if !errors.Is(err, ErrDBPing) {
		t.Error("db ping should have failed")
	}
}

// Given: test database is running on localhost:5432 with user/pass postgres.
func TestConnectWrongPort(t *testing.T) {
	dbFaultyConf := getTestDB()
	dbFaultyConf.Port = 84362

	_, err := Connect(dbFaultyConf)
	if !errors.Is(err, ErrDBPing) {
		t.Error("db ping should have failed")
	}
}
