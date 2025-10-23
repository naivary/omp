package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestReadyz(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{
			name: "readiness",
			code: http.StatusOK,
		},
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Errorf("new test sever: %s", err)
		t.FailNow()
	}
	endpoint, err := url.JoinPath(baseURL, "readyz")
	if err != nil {
		t.Errorf("URL join path: %s", err)
		t.FailNow()
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
