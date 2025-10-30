package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestLivez(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{
			name: "liveness",
			code: http.StatusOK,
		},
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("new test sever: %s", err)
	}
	endpoint, err := url.JoinPath(baseURL, "livez")
	if err != nil {
		t.Fatalf("URL join path: %s", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRequest[any](http.MethodGet, endpoint, nil)
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("do request: %s", err)
			}
			if res.StatusCode != tc.code {
				t.Fatalf("status code differ. Got: %d; Want: %d", res.StatusCode, tc.code)
			}
		})
	}
}
