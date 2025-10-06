package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/naivary/cnapi/probe"
)

var errProbeFailed = errors.New("probe failed")

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
	env := func(key string) string {
		switch key {
		case "HOST":
			return "127.0.0.1"
		case "PORT":
			return fmt.Sprintf("%d", port)
		default:
			return getenv(key)
		}
	}
	go run(ctx, args, env, stdin, stdout, stderr)

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
		return "", errProbeFailed
	}
	return baseURL, nil
}

// freePort returns a port which is probably useable.
func freePort() (int, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return -1, err
	}
	defer lis.Close()
	return lis.Addr().(*net.TCPAddr).Port, nil
}
