package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/naivary/omp/probe"
)

var errReadinessProbeFailed = errors.New("readiness probe failed")

// NewTestServer starts an HTTP server using the provided configuration parameters
// and ensures that it is fully initialized before accepting incoming connections.
//
// It wraps the `run` function to create an environment that closely mirrors
// production conditions. The function returns the server’s base URL, which can be
// used to send requests during testing.
//
// Typical usage in a test:
//
//	ctx := context.Background()
//	ctx, cancel := context.WithCancel(ctx)
//	t.Cleanup(cancel)
//
//	baseURL, err := NewTestServer(t, config)
//	require.NoError(t, err)
//	// The server is now ready to handle requests.
func NewTestServer(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (string, error) {
	port, err := freePort()
	if err != nil {
		return "", nil
	}
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	args = slices.Concat(args, []string{
		"omp",
		"-port", strconv.Itoa(port),
		"-pg.insert.testdata",
		"-pg.username", "postgres",
		"-pg.password", "postgres",
		"-pg.database", "omp",
		"-oidc.url", "http://127.0.0.1:8080",
		"-oidc.clientID", "omp-rest-api",
		"-oidc.clientSecret", getenv("OMP_OIDC_CLIENT_SECRET"),
	})
	go func() {
		err := run(ctx, args, getenv, stdin, stdout, stderr)
		if err != nil {
			panic(err)
		}
	}()
	// wait until ready
	readyzEndpoint, err := url.JoinPath(baseURL, "readyz")
	if err != nil {
		return "", err
	}
	r, err := http.NewRequest(http.MethodGet, readyzEndpoint, nil)
	if err != nil {
		return "", err
	}
	status, err := probe.DoHTTP(r, 5*time.Second)
	if err != nil {
		return "", err
	}
	if status == probe.Failed {
		return "", errReadinessProbeFailed
	}
	return baseURL, nil
}

// freePort returns a port which is probably useable. It's only PROBABLY useable
// because the listener is getting closed and their is a slight chance for
// another process to aquire that port.
func freePort() (int, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return -1, err
	}
	defer lis.Close()
	return lis.Addr().(*net.TCPAddr).Port, nil
}

func NewRequest[T any](method, url string, v T) *http.Request {
	r, err := newRequest(method, url, v)
	if err != nil {
		panic(err)
	}
	return r
}

func newRequest[T any](method, url string, v T) (*http.Request, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(data)
	return http.NewRequest(method, url, body)
}
