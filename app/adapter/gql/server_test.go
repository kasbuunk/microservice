package gql

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var (
	SvcPort        = 8089
	SvcGQLEndpoint = "/gql"
)

var conf = Config{
	Port:     SvcPort,
	Endpoint: SvcGQLEndpoint,
}
var serverURL = fmt.Sprintf("http://localhost:%v%v", conf.Port, conf.Endpoint)

func setupServer(t *testing.T) {
	svc, err := New(conf.Endpoint, nil)
	if err != nil {
		t.Error(err)
	}

	// Start server in separate process.
	go func() {
		err := svc.Serve(conf.Port)
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestGraphqlRequests(t *testing.T) {
	setupServer(t)

	testCases := []struct{ name, requestBody, expectedResponseBody string }{
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

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Send GET request to the endpoint.
			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, serverURL, strings.NewReader(testCase.requestBody))
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
			if actual != testCase.expectedResponseBody {
				t.Fatalf("expected: '%v'; actual: '%v'", testCase.expectedResponseBody, actual)
			}
		})
	}

}
