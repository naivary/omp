package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	clubv1 "github.com/naivary/omp/api/club/v1"
)

func TestCreateClub(t *testing.T) {
	tests := []struct {
		name string
		code int
		req  clubv1.CreateClubRequest
	}{
		{
			name: "valid request",
			code: http.StatusCreated,
			req: clubv1.CreateClubRequest{
				Email:    "info@omp.de",
				Password: "testpassword",
				Name:     "randomname",
				Location: "Berlin",
				Timezone: "Berlin/Europe",
			},
		},
	}
	ctx := context.Background()
	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("new test sever: %s", err)
	}
	endpoint, err := url.JoinPath(baseURL, "clubs")
	if err != nil {
		t.Fatalf("URL join path: %s", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRequest(http.MethodPost, endpoint, &tc.req)
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.code {
				io.Copy(os.Stdout, res.Body)
				t.Errorf("unexpected status code. Got: %d. Wanted: %d", res.StatusCode, tc.code)
			}
		})
	}
}
