package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"testing"

	clubv1 "github.com/naivary/omp/api/club/v1"
)

func getAllClubProfiles() ([]*clubv1.Profile, error) {
	return nil, nil
}

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
				t.Errorf("unexpected status code. Got: %d. Wanted: %d", res.StatusCode, tc.code)
			}
		})
	}
}

// func TestReadClub(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		code int
// 	}{
// 		{
// 			name: "valid request",
// 			code: http.StatusOK,
// 		},
// 	}
// 	ctx := context.Background()
// 	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
// 	if err != nil {
// 		t.Fatalf("new test sever: %s", err)
// 	}
// 	endpoint, err := url.JoinPath(baseURL, "clubs")
// 	if err != nil {
// 		t.Fatalf("URL join path: %s", err)
// 	}
//
// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			r := NewRequest[any](http.MethodPost, endpoint, nil)
// 			res, err := http.DefaultClient.Do(r)
// 			if err != nil {
// 				t.Fatalf("do request: %v", err)
// 			}
// 			defer res.Body.Close()
// 			if res.StatusCode != tc.code {
// 				t.Errorf("unexpected status code. Got: %d. Wanted: %d", res.StatusCode, tc.code)
// 			}
// 		})
// 	}
// }

func TestReadAllClubs(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{
			name: "valid request",
			code: http.StatusOK,
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
			r := NewRequest[any](http.MethodGet, endpoint, nil)
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.code {
				t.Errorf("unexpected status code. Got: %d. Wanted: %d", res.StatusCode, tc.code)
			}
			// there will be 100 club profiles provisioned beforehand
			numClubs := 100
			clubs := make([]*clubv1.Profile, 0, numClubs)
			err = json.NewDecoder(res.Body).Decode(&clubs)
			if err != nil {
				t.Fatalf("json decode: %v", err)
			}
			if len(clubs) < 100 {
				t.Errorf("number of club profiles not expected. Got: %d. Want: >=100", len(clubs))
			}
		})
	}
}
