// Package config loads global, immutable state that persists across the service's lifetime.
// It can include configuration populated by environment variables, static files and other
// sources of input that persist.
package config

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/kasbuunk/microservice/client/postmark"
	"github.com/kasbuunk/microservice/client/repository/storage"
	"github.com/kasbuunk/microservice/input/server"
)

const envPrefix = "svc"

// Config includes the data fields that the other microservice components need to set up.
type Config struct {
	Server   server.Config
	DB       storage.Config
	Postmark postmark.Config
}

// New takes its 'input' from environment variables and returns everything  the microservice
// needs to serve requests, including a port, endpoint and database configuration.
func New() (Config, error) {
	var conf Config
	err := envconfig.Process(envPrefix, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}
