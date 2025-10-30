package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	playerv1 "github.com/naivary/omp/api/player/v1"
)

func TestCreatePlayer(t *testing.T) {
	tests := []struct {
		name string
		code int
		req  playerv1.CreatePlayerRequest
	}{}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Errorf("new test sever: %s", err)
		t.FailNow()
	}
	endpoint, err := url.JoinPath(baseURL, "players")
	if err != nil {
		t.Errorf("URL join path: %s", err)
		t.FailNow()
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRequest(http.MethodPost, endpoint, tc.req)
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("send request: %v", err)
			}
			defer res.Body.Close()
			var errMsg bytes.Buffer
			io.Copy(&errMsg, res.Body)
			if res.StatusCode != tc.code {
				t.Logf("err msg: %s", errMsg.String())
				t.Fatalf("status code differ. Got: %d; Want: %d", res.StatusCode, tc.code)
			}
		})
	}
}
