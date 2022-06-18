package server

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/kasbuunk/microservice/test"
)

var conf = Config{
	Port:        Port(test.SvcPort),
	GQLEndpoint: GQLEndpoint(test.SvcGQLEndpoint),
}
var serverURL = fmt.Sprintf("http://localhost:%v%v", conf.Port, conf.GQLEndpoint)

func setupServer(t *testing.T) {
	server, err := New(conf.GQLEndpoint, nil)
	if err != nil {
		t.Error(err)
	}

	// Start server in separate process.
	go func() {
		err := server.Serve(conf.Port)
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestGraphqlRequests(t *testing.T) {
	setupServer(t)

	cases := []struct{ name, input, expected string }{
		{
			"EmptyRequest",
			"{}",
			"{\"errors\":[{\"message\":\"no operation provided\",\"extensions\":{\"code\":\"GRAPHQL_VALIDATION_FAILED\"}}],\"data\":null}",
		},
		{
			"GetUsers",
			"{\"query\":\"{ users { id email }}\"}",
			"{\"errors\":[{\"message\":\"internal system error\",\"path\":[\"users\"]}],\"data\":null}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Send GET request to the endpoint.
			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, serverURL, strings.NewReader(tc.input))
			if err != nil {
				t.Error(err)
			}
			req.Header.Add("Accept", "*/*")
			req.Header.Add("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)

			if err != nil {
				t.Error(err)
			}

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Error(err, "closing response body")
				}
			}(resp.Body)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err, "reading response body")
			}

			actual := string(body)
			if actual != tc.expected {
				t.Fatalf("expected: '%v'; actual: '%v'", tc.expected, actual)
			}
		})
	}

}
