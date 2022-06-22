// Package storage is the interface for a server's request handler to store and retrieve entities, following the
// domain's business logic.
package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq" // sql driver
)

const SQLDriver = "postgres"

var ErrDBConnection = errors.New("connecting to database")
var ErrDBPing = errors.New("pinging database")

type Config struct {
	Host string
	Port int
	Name string
	User string
	Pass string
}

func Connect(conf Config) (*sql.DB, error) {
	var db *sql.DB

	connection := fmt.Sprintf(
		"host=%s port=%d dbname=%s  user=%s  password=%s sslmode=disable",
		conf.Host, conf.Port, conf.Name, conf.User, conf.Pass)

	db, err := sql.Open(SQLDriver, connection)
	if err != nil {
		return db, ErrDBConnection
	}

	err = db.Ping()
	if err != nil {
		return nil, ErrDBPing
	}

	return db, nil
}
