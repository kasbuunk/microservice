package config

import (
	"strconv"
	"testing"
)

const ExpectedGotFormat string = "expected: '%v'; got: '%v'"

func TestNewConfig(t *testing.T) {
	port := 16543
	gqlEndpoint := "/my-gql-endpoint"
	dbHost := "thishost.me"
	dbPort := 2347
	dbName := "testname"
	dbUser := "thisuser"
	dbPass := "dontlook"

	// Set or override environment variables.
	t.Setenv("SVC_SERVER_PORT", strconv.Itoa(port))
	t.Setenv("SVC_SERVER_GQLENDPOINT", gqlEndpoint)
	t.Setenv("SVC_DB_HOST", dbHost)
	t.Setenv("SVC_DB_PORT", strconv.Itoa(dbPort))
	t.Setenv("SVC_DB_NAME", dbName)
	t.Setenv("SVC_DB_USER", dbUser)
	t.Setenv("SVC_DB_PASS", dbPass)

	conf, err := New()
	if err != nil {
		t.Error("getting config", err)
	}

	if conf.Server.Port != port {
		t.Errorf(ExpectedGotFormat, port, conf.Server.Port)
	}
	if conf.Server.GQLEndpoint != gqlEndpoint {
		t.Errorf(ExpectedGotFormat, gqlEndpoint, conf.Server.GQLEndpoint)
	}
	if conf.DB.Host != dbHost {
		t.Errorf(ExpectedGotFormat, dbHost, conf.DB.Host)
	}
	if conf.DB.Port != dbPort {
		t.Errorf(ExpectedGotFormat, dbPort, conf.DB.Port)
	}
	if conf.DB.Name != dbName {
		t.Errorf(ExpectedGotFormat, dbName, conf.DB.Name)
	}
	if conf.DB.User != dbUser {
		t.Errorf(ExpectedGotFormat, dbUser, conf.DB.User)
	}
	if conf.DB.Pass != dbPass {
		t.Errorf(ExpectedGotFormat, dbPass, conf.DB.Pass)
	}
}
