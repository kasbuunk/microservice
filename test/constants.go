// Package test provides constants and data to be used in other tests across the service.
package test

var (
	SvcPort        = 8089
	SvcGQLEndpoint = "/gql"
	DBHost         = "localhost"
	DBPort         = 5432
	DBName         = "auth_test"
	DBUser         = "postgres"
	DBPass         = "postgres"

	ExpectedGot = "expected '%v'; got '%v'"
)
