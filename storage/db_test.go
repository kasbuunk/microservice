package storage

import (
	"errors"
	"os"
	"testing"
)

func getTestDB() Config {
	host := os.Getenv("SVC_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	return Config{
		Host: host,
		Port: 5432,
		Name: "user_test",
		User: "postgres",
		Pass: "postgres",
	}
}

// Given: test database is running on localhost:5432 with user/pass postgres.
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
