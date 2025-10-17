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
	baseURL, err := NewTestServer(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
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
			r, err := http.NewRequest(http.MethodGet, endpoint, nil)
			if err != nil {
				t.Errorf("new request: %s", err)
				t.FailNow()
			}
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Errorf("do request: %s", err)
				t.FailNow()
			}
			if res.StatusCode != tc.code {
				t.Errorf("status code differ. Got: %d; Want: %d", res.StatusCode, tc.code)
				t.FailNow()
			}
		})
	}
}
