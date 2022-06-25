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

	serverToken := "myservertoken"
	accountToken := "myaccounttoken"

	// Set or override environment variables.
	t.Setenv("SVC_GQLSERVER_PORT", strconv.Itoa(port))
	t.Setenv("SVC_GQLSERVER_ENDPOINT", gqlEndpoint)

	t.Setenv("SVC_DB_HOST", dbHost)
	t.Setenv("SVC_DB_PORT", strconv.Itoa(dbPort))
	t.Setenv("SVC_DB_NAME", dbName)
	t.Setenv("SVC_DB_USER", dbUser)
	t.Setenv("SVC_DB_PASS", dbPass)

	t.Setenv("SVC_POSTMARK_SERVERTOKEN", serverToken)
	t.Setenv("SVC_POSTMARK_ACCOUNTTOKEN", accountToken)

	conf, err := New()
	if err != nil {
		t.Error("getting config", err)
	}

	if int(conf.GQLServer.Port) != port {
		t.Errorf(ExpectedGotFormat, port, conf.GQLServer.Port)
	}
	if string(conf.GQLServer.Endpoint) != gqlEndpoint {
		t.Errorf(ExpectedGotFormat, gqlEndpoint, conf.GQLServer.Endpoint)
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
